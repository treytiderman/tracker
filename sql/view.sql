-- all trackers
-- SELECT * FROM tracker;

-- -- all field
-- SELECT * FROM field;

-- -- all number
-- SELECT * FROM number;

-- -- all option
-- SELECT * FROM option;



-- all trackers and their fields
SELECT * FROM tracker
LEFT JOIN field USING (tracker_id);

-- all fields and their types
SELECT * FROM field
LEFT JOIN number USING (field_id)
LEFT JOIN option USING (field_id);

-- all trackers and their fields and their types
SELECT
    tracker_id,
    tracker_name,
    tracker_notes,

    IFNULL(field_id, 0) AS field_id,
    IFNULL(field_type, "") AS field_type,
    IFNULL(field_name, "") AS field_name,
    IFNULL(field_notes, "") AS field_notes,

    IFNULL(number_id, 0) AS number_id,
    IFNULL(decimal_places, 0) AS decimal_places,

    IFNULL(option_id, 0) AS option_id,
    IFNULL(option_value, 0) AS option_value,
    IFNULL(option_name, "") AS option_name
FROM (
    SELECT * FROM tracker
    LEFT JOIN field USING (tracker_id)
    -- WHERE tracker_name = "Brush Teeth"
    -- WHERE tracker_id = 2
) AS tf
LEFT JOIN number AS n USING (field_id)
LEFT JOIN option AS o USING (field_id)
ORDER BY tf.tracker_id, tf.field_id, n.number_id, o.option_id;



-- simple entries
SELECT * FROM entry;

-- simple entries with logs
SELECT * FROM entry
LEFT JOIN log USING (entry_id);

-- entries with logs and preset column
SELECT 
    entry.entry_id,
    entry.timestamp,
    entry.entry_notes,

    IFNULL(log_id, 0) AS log_id,
    IFNULL(log_value, 0) AS log_value,

    IFNULL(field.field_id, 0) AS field_id,
    IFNULL(field_type, "") AS field_type,

    IFNULL(decimal_places, 0) AS decimal_places,

    IFNULL(option_value, 0) AS option_value,
    IFNULL(option_name, "") AS option_name,

    IFNULL((CASE WHEN field.field_type == "number" THEN 
        printf("%.2f", log.log_value / power(10, number.decimal_places))
     ELSE 
        option.option_name
     END), "") AS present
FROM entry
LEFT JOIN log USING (entry_id)
LEFT JOIN field USING (field_id)
LEFT JOIN number USING (field_id)
LEFT JOIN option ON log.log_value = option.option_value
WHERE entry.tracker_id = 8
ORDER BY entry.entry_id, log.log_id;
