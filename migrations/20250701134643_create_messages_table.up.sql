CREATE TABLE messages
(
    id          SERIAL PRIMARY KEY,
    sender_id   INT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    room_id     INT REFERENCES rooms (id) ON DELETE CASCADE,
    receiver_id INT REFERENCES users (id) ON DELETE CASCADE,
    content     TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        (room_id IS NOT NULL AND receiver_id IS NULL) OR
        (room_id IS NULL AND receiver_id IS NOT NULL)
        )
);