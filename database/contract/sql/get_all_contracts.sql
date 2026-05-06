select id, number, subscriber_id, object_id, sign_date, created_at, updated_at
from contracts
order by id
limit $1 offset $2;
