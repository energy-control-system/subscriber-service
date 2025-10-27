select id,
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
       updated_at
from subscribers
where id = $1;
