-- Update Tracker
UPDATE tracker
SET tracker_notes = "Every time I take the dogs out",
    tracker_name = "Walk Dogs"
WHERE tracker_id = 1;

-- Update Entry
UPDATE entry
SET entry_notes = "Dog found a turtle dove",
    timestamp = "2730-01-09 14:16:00"
WHERE entry_id = 2;

-- Update Log
UPDATE log SET log_value = 42069 WHERE log_id = 3;

-- Update Field
-- Changing field_type is not allowed
UPDATE field
SET field_name = "Purchases",
    field_notes = "Purchases of different types"
WHERE field_id = 1;

-- Update number
-- Decreasing decimal_places will lose that decimal place data
UPDATE log
SET log_value = log_value * power(10, (3) - (SELECT decimal_places FROM number WHERE field_id = 1))
WHERE field_id = 1;

UPDATE number
SET decimal_places = (3)
WHERE field_id = 1;

-- Update option
-- Also updates any logs that contain the value of the option updated
UPDATE log
SET log_value = 4
FROM (SELECT field_id, option_value FROM option WHERE option_id = 3) AS o
WHERE log.field_id = o.field_id AND log.log_value = o.option_value;

UPDATE option
SET option_value = 4,
    option_name = "Mastercard"
WHERE option_id = 3;

-- Add option
INSERT INTO option (field_id, option_value, option_name)
VALUES (2, 5, "Cash");

-- Add option - Test data to delete
INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Gold");
INSERT INTO log (entry_id, field_id, log_value) VALUES (8, 1, 50000);
INSERT INTO log (entry_id, field_id, log_value) VALUES (8, 2, 5);

-- Delete option - WIP
-- Also delete any logs + entries that contain the value of the option deleted
DELETE FROM entry WHERE entry_id = (
    SELECT entry_id FROM log
    WHERE log.field_id = (SELECT field_id FROM option WHERE option_id = 4)
    AND log.log_value = (SELECT option_value FROM option WHERE option_id = 4)
);

DELETE FROM option WHERE option_id = (4);

-- Last SELECT
SELECT * FROM tracker
