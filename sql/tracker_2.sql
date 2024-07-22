-- Delete Tracker
DELETE FROM tracker WHERE tracker_name = "Money";

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Money", "Spending habits");

-- Add Data Fields to Tracker
INSERT INTO field (tracker_id, field_type, field_name, field_notes)
VALUES (2, "number", "Transactions", "Transactions of different types");
INSERT INTO number (field_id, decimal_places)
VALUES (1, 2);

INSERT INTO field (tracker_id, field_type, field_name, field_notes)
VALUES (2, "option", "Card", "What credit/debit card was used?");
INSERT INTO option (field_id, option_value, option_name)
VALUES (2, 1, "Visa"),
       (2, 2, "Discover"),
       (2, 3, "American Express");

-- Log values
INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Gum");
INSERT INTO log (entry_id, field_id, log_value) VALUES (5, 1, 1099);
INSERT INTO log (entry_id, field_id, log_value) VALUES (5, 2, 2);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Turtle");
INSERT INTO log (entry_id, field_id, log_value) VALUES (6, 1, 42050);
INSERT INTO log (entry_id, field_id, log_value) VALUES (6, 2, 3);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Coffee");
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 1, 799);
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 2, 1);

-- Select All entries
SELECT * FROM entry WHERE tracker_id = 2;
