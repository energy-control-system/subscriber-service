select id, number, subscriber_id, object_id, sign_date, created_at, updated_at
from contracts
where object_id = $1
order by sign_date
limit 1;
