-- Update Tracker
UPDATE tracker
SET tracker_notes = "Every time I take the dogs out",
    tracker_name = "Walk Dogs"
WHERE tracker_id = 1;

-- Update Entry
UPDATE entry
SET entry_notes = "Dog found a turtle dove",
    timestamp = "2030-01-09 14:16:00"
WHERE entry_id = 2;

-- Update Log
UPDATE log
SET log_value = 42069
WHERE entry_id = 6 AND field_id = 1;

-- Update Field
-- Changing field_type is not allowed
UPDATE field
SET field_name = "Purchases",
    field_notes = "Purchases of different types"
WHERE field_id = 1;

-- Update number
-- Decreasing decimal_places will lose log_value precision
UPDATE log
SET log_value = log_value * power(10, (3) - (SELECT decimal_places FROM number WHERE field_id = 1))
WHERE field_id = 1;

UPDATE number
SET decimal_places = (3)
WHERE field_id = 1;

-- Update option
UPDATE option
SET option_value = 4,
    option_name = "Mastercard"
WHERE option_id = 3;

-- Add option
-- INSERT INTO option (field_id, option_value, option_name)
-- VALUES (2, 5, "Cash");

-- Delete option
-- DELETE FROM log WHERE field_id = o.field_id AND log_value = o.option_value (
--     SELECT field_id, option_value FROM option WHERE option_id = 1
-- ) AS o;
-- DELETE FROM option WHERE option_id = 1;

-- DELETE FROM log WHERE field_id = 2 AND log_value = 2;


SELECT field_id, option_value FROM option WHERE option_id = 1
