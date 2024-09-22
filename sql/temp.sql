-- SELECT entry_id FROM log
-- WHERE log.field_id = (SELECT field_id FROM option WHERE option_id = 4)
-- AND log.log_value = (SELECT option_value FROM option WHERE option_id = 4);

DELETE FROM entry WHERE entry_id = (
    SELECT entry_id FROM log
    WHERE log.field_id = (SELECT field_id FROM option WHERE option_id = 4)
    AND log.log_value = (SELECT option_value FROM option WHERE option_id = 4)
);

DELETE FROM option WHERE option_id = (4);

SELECT * FROM entry;
