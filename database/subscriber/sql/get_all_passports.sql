select distinct on (subscriber_id) id, subscriber_id, series, number, issued_by, issue_date
from passports
order by subscriber_id, issue_date desc;
