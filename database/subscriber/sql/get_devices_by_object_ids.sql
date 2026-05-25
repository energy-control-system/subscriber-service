select id, object_id, type, number, place_type, place_description, created_at, updated_at
from devices
where object_id = any ($1)
order by id;
