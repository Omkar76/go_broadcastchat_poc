-- CREATE DATABASE chatapp;
-- use chatapp;
--  postgres queries
CREATE TABLE users (
    id SERIAL NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

-- INSERT INTO users (username, password) VALUES ('admin', 'admin');
-- INSERT INTO users (username, password) VALUES ('user', 'user');

CREATE TABLE messages (
    id SERIAL NOT NULL,
    message VARCHAR(255) NOT NULL,
    from INTEGER NOT NULL,
    to INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE friends (
    id SERIAL NOT NULL,
    user_id INTEGER NOT NULL,
    friend_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (friend_id) REFERENCES users (id)
);

CREATE TABLE groups (
    id SERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE group_members (
    id SERIAL NOT NULL,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);