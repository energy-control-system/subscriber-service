select id, address, have_automaton, created_at, updated_at
from objects
where id = $1;
