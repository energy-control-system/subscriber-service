update contracts
set number = :number,
    subscriber_id = :subscriber_id,
    object_id = :object_id,
    sign_date = :sign_date
where id = :id
returning id, number, subscriber_id, object_id, sign_date, created_at, updated_at;
