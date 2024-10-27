PRAGMA foreign_keys = ON;

-- Drop Tables
-----------------------------------------------------------------------------------------

DROP TABLE IF EXISTS content;
DROP TABLE IF EXISTS log;
DROP TABLE IF EXISTS entry;
DROP TABLE IF EXISTS option;
DROP TABLE IF EXISTS number;
DROP TABLE IF EXISTS field;
DROP TABLE IF EXISTS tracker;


-- Create Tables
-----------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS tracker (
    tracker_id INTEGER NOT NULL UNIQUE,
    tracker_name TEXT NOT NULL UNIQUE,
    tracker_notes TEXT NOT NULL DEFAULT "",
    PRIMARY KEY (tracker_id)
);

CREATE TABLE IF NOT EXISTS field (
    field_id INTEGER NOT NULL UNIQUE,
    tracker_id INTEGER NOT NULL,
    field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',
    field_name TEXT NOT NULL,
    field_notes TEXT NOT NULL DEFAULT "",
    UNIQUE(tracker_id, field_name),
    PRIMARY KEY (field_id),
    FOREIGN KEY (tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS number (
    number_id INTEGER NOT NULL UNIQUE,
    field_id INTEGER NOT NULL,
    decimal_places INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (number_id),
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS option (
    option_id INTEGER NOT NULL UNIQUE,
    field_id INTEGER NOT NULL,
    option_value INTEGER NOT NULL DEFAULT 0,
    option_name TEXT NOT NULL DEFAULT "value",
    UNIQUE(field_id, option_name),
    PRIMARY KEY (option_id),
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS entry (
    entry_id INTEGER NOT NULL UNIQUE,
    tracker_id INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    entry_notes TEXT NOT NULL DEFAULT "",
    PRIMARY KEY (entry_id),
    FOREIGN KEY (tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS log (
    log_id INTEGER NOT NULL UNIQUE,
    entry_id INTEGER NOT NULL,
    field_id INTEGER NOT NULL,
    log_value INTEGER NOT NULL,

    PRIMARY KEY (log_id),
    FOREIGN KEY (entry_id) REFERENCES entry (entry_id) ON DELETE CASCADE,
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS content (
    content_id INTEGER NOT NULL UNIQUE,
    entry_id INTEGER NOT NULL,
    content_path TEXT NOT NULL DEFAULT "",
    PRIMARY KEY (content_id),
    FOREIGN KEY (entry_id) REFERENCES entry (entry_id)
);


-- Tracker 1 - Journal
-----------------------------------------------------------------------------------------

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Journal", "Searchable journal expecting markdown syntax");

-- Log values
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog found a turtle");
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog learned to fly");
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog ran away");
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "No dog still walk");
INSERT INTO entry (tracker_id, entry_notes) VALUES (1,
"# Basic Markdown

# Heading 1
## Heading 2
### Heading 3

Text **bold** *italicized* `code` text

Horizontal Rule

---

> Blockquote

Ordered List

1. First item
2. Second item
3. Third item

Unordered List

- First item
- Second item
- Third item

[Google Link](https://www.google.com)

![Image Link](image.jpg)
");
INSERT INTO entry (tracker_id, entry_notes) VALUES (1,
'# Extended Markdown

~~Strikethrough~~

Table

| Syntax | Description |
| ----------- | ----------- |
| Header | Title |
| Paragraph | Text |

Code Block

```json
{
    "firstName": "John",
    "lastName": "Smith",
    "age": 25
}
```

```go
func Routes_pages(db *sql.DB) {
    page_Trackers(db)
}
```

Task List

- [x] Task 1
- [ ] Task 2
- [ ] Task 3

');
INSERT INTO content (entry_id, content_path) VALUES (1, "/content/2024-10-26_08:17:45.png");
INSERT INTO content (entry_id, content_path) VALUES (1, "/content/2024-10-26_09:34:11.png");

-- Tracker 2 - Money
-- --------------------------------------------------------------------------------------

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
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 1, 1099);
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 2, 2);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Turtle");
INSERT INTO log (entry_id, field_id, log_value) VALUES (8, 1, 42050);
INSERT INTO log (entry_id, field_id, log_value) VALUES (8, 2, 3);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Coffee");
INSERT INTO log (entry_id, field_id, log_value) VALUES (9, 1, 799);
INSERT INTO log (entry_id, field_id, log_value) VALUES (9, 2, 1);
