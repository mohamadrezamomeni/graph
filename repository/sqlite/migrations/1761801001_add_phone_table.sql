
-- +migrate Up
CREATE TABLE phones (
    phone VARCHAR(20) PRIMARY KEY,
    contact_id INTEGER NOT NULL,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE `phones`;


