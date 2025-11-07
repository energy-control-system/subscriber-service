insert
into devices (object_id, type, number, place_type, place_description)
select o.id, :type, :number, :place_type, :place_description
from objects o
where address = :object_address
on conflict (number) do update
    set type              = excluded.type,
        place_type        = excluded.place_type,
        place_description = excluded.place_description;
