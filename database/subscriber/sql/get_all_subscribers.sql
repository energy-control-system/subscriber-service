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
from subscribers s
where ($1 = ''
    or s.surname ilike '%' || $1 || '%'
    or s.name ilike '%' || $1 || '%'
    or s.patronymic ilike '%' || $1 || '%'
    or s.account_number ilike '%' || $1 || '%'
    or s.phone_number ilike '%' || $1 || '%'
    or exists (select 1
               from contracts c
                        join objects o on o.id = c.object_id
               where c.subscriber_id = s.id
                 and o.address ilike '%' || $1 || '%'))
order by id
limit $2 offset $3;
