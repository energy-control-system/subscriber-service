insert into contracts (number, subscriber_id, object_id, sign_date)
values (:number, :subscriber_id, :object_id, :sign_date)
returning id, number, subscriber_id, object_id, sign_date, created_at, updated_at;
