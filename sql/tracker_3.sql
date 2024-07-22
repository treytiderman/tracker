DROP TABLE IF EXISTS tracker_3;
DELETE FROM tracker WHERE tracker_name = "Lift";

INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Lift", "lifting habits");

INSERT INTO field 
    (tracker_id, field_type, field_name)
VALUES 
    (3, "option", "Exercise"),
    (3, "number", "Weight"),
    (3, "number", "Reps");

INSERT INTO option
    (field_id, option_value, option_name)
VALUES 
    (2, 1, "Bench"),
    (2, 2, "Squat"),
    (2, 3, "Deadlift");

INSERT INTO number (field_id)
VALUES (3), (4);

SELECT tracker_id FROM tracker WHERE tracker_name = "Lift";

CREATE TABLE IF NOT EXISTS tracker_3 (
    id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL DEFAULT 3,

    -- timestamp when this was set
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- markdown formatted notes
    notes TEXT NOT NULL DEFAULT "",

    -- custom fields from the field table
    "Exercise" INTEGER NOT NULL DEFAULT 0,
    "Weight" INTEGER NOT NULL DEFAULT 0,
    "Reps" INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(id),
    FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id)
);

INSERT INTO tracker_3
    (notes, "Exercise", "Weight", "Reps")
VALUES
    ("", 1, 135, 5),
    ("finally two plates!!", 2, 225, 5),
    ("", 1, 135, 8),
    ("broke back after first rep :(", 2, 225, 1);

SELECT * FROM tracker_3;