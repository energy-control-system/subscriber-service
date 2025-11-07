with subscriber as (select id
                    from subscribers
                    where account_number = :subscriber_account_number),
     object as (select id
                from objects
                where address = :object_address)
insert
into contracts (number, subscriber_id, object_id, sign_date)
select :number, s.id, o.id, :sign_date
from subscriber s,
     object o
on conflict (number) do update
    set sign_date = excluded.sign_date;
