update subscribers
set account_number = :account_number,
    surname = :surname,
    name = :name,
    patronymic = :patronymic,
    phone_number = :phone_number,
    email = :email,
    inn = :inn,
    birth_date = :birth_date,
    status = case when :status = 0 then status else :status end
where id = :id
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
