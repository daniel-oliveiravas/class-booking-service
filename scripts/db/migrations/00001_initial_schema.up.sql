CREATE TABLE IF NOT EXISTS members
(
    id         TEXT      NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name       TEXT      NOT NULL
);

CREATE TABLE IF NOT EXISTS classes
(
    id         TEXT      NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name       TEXT      NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date   TIMESTAMP NOT NULL,
    capacity   int       NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings
(
    id         TEXT      NOT NULL PRIMARY KEY,
    member_id  TEXT      NOT NULL,
    class_id   TEXT      NOT NULL,
    class_date DATE      NOT NULL,
    booked_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (member_id) REFERENCES members (id),
    FOREIGN KEY (class_id) REFERENCES classes (id)
);
