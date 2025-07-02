CREATE TABLE rooms
(
    id         SERIAL PRIMARY KEY,
    name       varchar(100) NOT NULL UNIQUE,
    is_private BOOLEAN      NOT NULL DEFAULT FALSE,
    created_by INT          REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
)