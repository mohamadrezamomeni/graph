-- +migrate Up
CREATE TABLE contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name VARCHAR(32),
    last_name VARCHAR(32)
);

-- +migrate Down
DROP TABLE contacts;