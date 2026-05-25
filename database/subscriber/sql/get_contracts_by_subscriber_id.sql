select id, number, subscriber_id, object_id, sign_date, created_at, updated_at
from contracts
where subscriber_id = $1
order by id;
