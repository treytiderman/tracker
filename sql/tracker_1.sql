-- Delete Tracker
DELETE FROM tracker WHERE tracker_name = "Walk Dog";

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Walk Dog", "idk yet");

-- Update Tracker Notes
UPDATE tracker
SET tracker_notes = "Take the dog out"
WHERE tracker_name = "Walk Dog";

UPDATE tracker
SET tracker_notes = "Every time I take the dog out"
WHERE tracker_id = 1;

-- Log values
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog found a turtle");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog learned to fly");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog ran away");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "No dog still walk");

-- Select All entries
SELECT * FROM entry;
