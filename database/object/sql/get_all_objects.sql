select id, address, have_automaton, created_at, updated_at
from objects
order by id
limit $1 offset $2;
