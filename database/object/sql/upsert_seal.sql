insert
into seals (device_id, number, place)
select d.id, :number, :place
from devices d
where d.number = :device_number
on conflict (number) do update
    set place = excluded.place;
