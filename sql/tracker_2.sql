DROP TABLE IF EXISTS tracker_2;
DELETE FROM tracker WHERE tracker_name = "money";

INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("money", "spending habits");

INSERT INTO field (tracker_id, field_type, field_name)
VALUES (2, "number", "transactions");

INSERT INTO number (field_id, decimal_places)
VALUES (1, 2);

SELECT tracker_id FROM tracker WHERE tracker_name = "money";

CREATE TABLE IF NOT EXISTS tracker_2 (
    id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL DEFAULT 2,

    -- timestamp when this was set
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- markdown formated notes
    notes TEXT NOT NULL DEFAULT "",

    -- custom fields from the field table
    transactions INT NOT NULL DEFAULT 0,

    PRIMARY KEY(id),
    FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id)
);

INSERT INTO tracker_2
    (notes, transactions)
VALUES
    ("found a $100 bill", 10000),
    ("lost a bet", -10000),
    ("bought some gum", -200),
    ("paycheck", 100000);

SELECT * FROM tracker_2;