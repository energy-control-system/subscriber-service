select distinct on (subscriber_id) id, subscriber_id, series, number, issued_by, issue_date
from passports
where subscriber_id = any ($1)
order by subscriber_id, issue_date desc;
