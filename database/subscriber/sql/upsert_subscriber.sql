with subscriber as (
    insert into subscribers (account_number, surname, name, patronymic, phone_number, email, inn, birth_date)
        values (:account_number, :surname, :name, :patronymic, :phone_number, :email, :inn, :birth_date)
        on conflict (account_number) do update
            set surname = excluded.surname,
                name = excluded.name,
                patronymic = excluded.patronymic,
                phone_number = excluded.phone_number,
                email = excluded.email,
                inn = excluded.inn,
                birth_date = excluded.birth_date
        returning id)
insert
into passports (subscriber_id, series, number, issued_by, issue_date)
select s.id, :passport_series, :passport_number, :passport_issued_by, :passport_issue_date
from subscriber s
on conflict (series, number) do update
    set issued_by  = excluded.issued_by,
        issue_date = excluded.issue_date;
