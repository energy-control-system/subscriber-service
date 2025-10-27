insert into objects (address, have_automaton)
values (:address, :have_automaton)
returning id, address, have_automaton, created_at, updated_at;
