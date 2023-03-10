-- +migrate Up
CREATE TABLE postal_api_responses
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    api_name         TEXT    NOT NULL,
    tracking_number  TEXT    NOT NULL,
    first_fetched_at INTEGER NOT NULL,
    last_fetched_at  INTEGER NOT NULL,
    response_body    TEXT    NOT NULL,
    status           TEXT    NOT NULL
);


-- +migrate Down
DROP TABLE postal_api_responses;
