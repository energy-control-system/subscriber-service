update passports
set series = :passport_series,
    number = :passport_number,
    issued_by = :passport_issued_by,
    issue_date = :passport_issue_date
where subscriber_id = :id
returning id, subscriber_id, series, number, issued_by, issue_date;
