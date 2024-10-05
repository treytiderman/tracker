-- Delete Tracker
DELETE FROM tracker WHERE tracker_name = "Walk Dog";

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Walk Fish", "idk yet");

-- Update Tracker Notes
UPDATE tracker
SET tracker_notes = "Every time I take the dog out",
    tracker_name = "Walk Dog"
WHERE tracker_id = 3;

-- Add Entry
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog found a turtle");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog learned to fly");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog ran away");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "No dog still walk");

-- Update Entry
UPDATE entry
SET entry_notes = "Dog found a turtle dove",
    timestamp = "2030-09-09 14:16:00"
WHERE entry_id = 1;

-- Select All entries
SELECT * FROM entry;
