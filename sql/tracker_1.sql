DROP TABLE IF EXISTS tracker_1;
DELETE FROM tracker WHERE tracker_name = "walk dog";

INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("walk dog", "each time i take the dog out");

SELECT tracker_id FROM tracker WHERE tracker_name = "walk dog";

CREATE TABLE IF NOT EXISTS tracker_1 (
    id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL DEFAULT 1,

    -- timestamp when this was set
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- markdown formated notes
    notes TEXT NOT NULL DEFAULT "",

    PRIMARY KEY(id),
    FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

INSERT INTO tracker_1
    (notes)
VALUES
    ("dog found a turtle"),
    ("dog learned to fly"),
    ("dog ran away"),
    ("no dog still walk");

SELECT * FROM tracker_1;