-- all trackers
SELECT * FROM tracker;

-- all field
SELECT * FROM field;

-- all number
SELECT * FROM number;

-- all option
SELECT * FROM option;

-- all trackers and their fields
SELECT * FROM tracker
LEFT JOIN field USING (tracker_id);

-- all fields and their types
SELECT * FROM field
LEFT JOIN number USING (field_id)
LEFT JOIN option USING (field_id);

-- all trackers and their fields and their types
SELECT * FROM
(
    SELECT * FROM tracker
    LEFT JOIN field USING (tracker_id)
) tf
LEFT JOIN number n USING (field_id)
LEFT JOIN option o USING (field_id)
ORDER BY tf.field_id, o.option_id;
