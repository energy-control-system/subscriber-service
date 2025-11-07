insert into objects (address, have_automaton)
values (:address, :have_automaton)
on conflict (address) do update
    set have_automaton = excluded.have_automaton;
