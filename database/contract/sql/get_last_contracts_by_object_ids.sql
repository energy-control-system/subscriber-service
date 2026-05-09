select distinct on (object_id) id,
                               number,
                               subscriber_id,
                               object_id,
                               sign_date,
                               created_at,
                               updated_at
from contracts
where object_id in (?)
order by object_id, sign_date desc, id desc;
