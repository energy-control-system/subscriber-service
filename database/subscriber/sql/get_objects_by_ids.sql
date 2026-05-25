select id, address, have_automaton, created_at, updated_at
from objects
where id = any ($1)
order by id;
