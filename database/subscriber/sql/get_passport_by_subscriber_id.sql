select id, subscriber_id, series, number, issued_by, issue_date
from passports
where subscriber_id = $1
order by issue_date
limit 1;
