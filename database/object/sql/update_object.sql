update objects
set address = :address,
    have_automaton = :have_automaton
where id = :id
returning id, address, have_automaton, created_at, updated_at;
