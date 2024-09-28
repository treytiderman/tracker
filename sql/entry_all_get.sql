-- entries with logs and preset column
SELECT
    entry.entry_id,
    entry.tracker_id,
    entry.timestamp,
    entry.entry_notes,
    IFNULL(log.log_id, 0) AS log_id,
    IFNULL(log.log_value, 0) AS log_value,
    IFNULL(field.field_id, 0) AS field_id,
    IFNULL(field.field_type, "") AS field_type,
    IFNULL(field.field_name, "") AS field_name,
    IFNULL(number.decimal_places, 0) AS decimal_places,
    IFNULL(option.option_value, 0) AS option_value,
    IFNULL(option.option_name, "") AS option_name,
    IFNULL((CASE WHEN field.field_type == "number" THEN
        printf(("%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
    ELSE
        option.option_name
    END), "") AS present
FROM entry
LEFT JOIN log USING (entry_id)
LEFT JOIN field USING (field_id)
LEFT JOIN number USING (field_id)
LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
WHERE entry.tracker_id = 2
ORDER BY entry.entry_id, field.field_id;
