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
where ($1 = '' or s.surname ilike '%' || $1 || '%')
  and ($2 = '' or s.name ilike '%' || $2 || '%')
  and ($3 = '' or s.patronymic ilike '%' || $3 || '%')
  and ($4 = '' or s.account_number ilike '%' || $4 || '%')
  and ($5 = '' or s.phone_number ilike '%' || $5 || '%')
  and ($6 = '' or exists (select 1
                          from contracts c
                                   join objects o on o.id = c.object_id
                          where c.subscriber_id = s.id
                            and o.address ilike '%' || $6 || '%'))
order by id
limit $7 offset $8;
