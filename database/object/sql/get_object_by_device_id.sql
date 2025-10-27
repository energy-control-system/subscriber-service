select o.id, o.address, o.have_automaton, o.created_at, o.updated_at
from objects o
         join devices d on o.id = d.object_id
where d.id = $1;
