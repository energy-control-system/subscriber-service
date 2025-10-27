insert into passports (subscriber_id, series, number, issued_by, issue_date)
values (:subscriber_id, :series, :number, :issued_by, :issue_date)
returning id, subscriber_id, series, number, issued_by, issue_date;
