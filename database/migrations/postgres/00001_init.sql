-- +goose Up
create table if not exists subscriber_statuses
(
    id   int primary key generated always as identity,
    name text not null
);

insert into subscriber_statuses (name)
values ('Active'),
       ('Violator'),
       ('Archived');

create table if not exists subscribers
(
    id             int primary key generated always as identity,
    account_number text        not null unique,                        -- Лицевой счет
    surname        text        not null,
    name           text        not null,
    patronymic     text        not null,
    phone_number   text        not null check ( phone_number ~ '^(\+7|8)\d{10}$' ),
    email          text        not null check ( email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$' ),
    inn            text        not null check ( inn ~ '^\d{10,12}$' ), -- ИНН
    birth_date     date        not null check ( '1900-01-01' <= birth_date and birth_date < current_date ),
    status         int         not null default 1 references subscriber_statuses (id) on delete restrict,
    created_at     timestamptz not null default now(),
    updated_at     timestamptz not null default now()
);

create table if not exists passports
(
    id            int primary key generated always as identity,
    subscriber_id int  not null references subscribers (id) on delete cascade,
    series        text not null check ( series ~ '^\d{4}$' ), -- Серия
    number        text not null check ( number ~ '^\d{6}$' ), -- Номер
    issued_by     text not null,                              -- Кем выдан
    issue_date    date not null,                              -- Когда выдан
    unique (series, number)
);

create table if not exists objects
(
    id             int primary key generated always as identity,
    address        text        not null unique,
    have_automaton bool        not null, -- Наличие коммутационного (вводного) аппарата
    created_at     timestamptz not null default now(),
    updated_at     timestamptz not null default now()
);

create table if not exists device_place_types
(
    id   int primary key generated always as identity,
    name text not null
);

insert into device_place_types (name)
values ('Other'),
       ('Flat'),
       ('StairLanding');

create table if not exists devices
(
    id                int primary key generated always as identity,
    object_id         int         not null references objects (id) on delete cascade,
    type              text        not null,
    number            text        not null unique,
    place_type        int         not null references device_place_types (id) on delete restrict,
    place_description text        not null, -- Место установки прибора учета
    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now()
);

create table if not exists seals
(
    id         int primary key generated always as identity,
    device_id  int         not null references devices (id) on delete cascade,
    number     text        not null unique,
    place      text        not null, -- Место установки
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists contracts
(
    id            int primary key generated always as identity,
    number        text        not null unique,
    subscriber_id int         not null references subscribers (id) on delete restrict,
    object_id     int         not null references objects (id) on delete restrict,
    sign_date     date        not null check ( sign_date <= current_date ), -- Дата подписания
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now()
);

create index if not exists idx_subscribers_status on subscribers (status);
create index if not exists idx_passports_subscriber on passports (subscriber_id);
create index if not exists idx_devices_object on devices (object_id);
create index if not exists idx_devices_place_type on devices (place_type);
create index if not exists idx_seals_device on seals (device_id);
create index if not exists idx_contracts_subscriber on contracts (subscriber_id);
create index if not exists idx_contracts_object on contracts (object_id);

-- +goose StatementBegin
create or replace function update_updated_at_column()
    returns trigger as
$$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
do
$$
    declare
        table_name text;
    begin
        for table_name in
            select t.table_name
            from information_schema.tables t
                     join information_schema.columns c on t.table_name = c.table_name
            where t.table_schema = 'public'
              and c.column_name = 'updated_at'
            loop
                execute format('create trigger trg_%I_updated_at
                       before update on %I
                       for each row execute function update_updated_at_column()',
                               table_name, table_name);
            end loop;
    end
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
do
$$
    declare
        table_name text;
    begin
        for table_name in
            select t.table_name
            from information_schema.tables t
                     join information_schema.columns c on t.table_name = c.table_name
            where t.table_schema = 'public'
              and c.column_name = 'updated_at'
            loop
                execute format('drop trigger if exists trg_%I_updated_at on %I', table_name, table_name);
            end loop;
    end
$$;
-- +goose StatementEnd

drop function if exists update_updated_at_column();
drop table if exists contracts;
drop table if exists seals;
drop table if exists devices;
drop table if exists device_place_types;
drop table if exists objects;
drop table if exists passports;
drop table if exists subscribers;
drop table if exists subscriber_statuses;
