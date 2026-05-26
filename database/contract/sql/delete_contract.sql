delete from contracts
where id = $1
returning id, number, subscriber_id, object_id, sign_date, created_at, updated_at;
