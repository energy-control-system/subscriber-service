select s.id, s.device_id, s.number, s.place, s.created_at, s.updated_at
from seals s
         join devices d on d.id = s.device_id
where d.object_id = $1;
