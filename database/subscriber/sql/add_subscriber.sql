insert into subscribers (account_number, surname, name, patronymic, phone_number, email, inn, birth_date, status)
values (:account_number, :surname, :name, :patronymic, :phone_number, :email, :inn, :birth_date, :status)
returning id,
    account_number,
    surname,
    name,
    patronymic,
    phone_number,
    email,
    inn,
    birth_date,
    status,
    created_at,
    updated_at;
