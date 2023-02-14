CREATE TABLE users
(
    id            INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(255)    NOT NULL,
    username      VARCHAR(255)    NOT NULL UNIQUE,
    password_hash VARCHAR(255)    NOT NULL
);

CREATE TABLE todo_lists
(
    id          INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title       VARCHAR(255)    NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE users_lists
(
    id      INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT             NOT NULL,
    list_id INT             NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (list_id) REFERENCES todo_lists (id)
);

CREATE TABLE todo_items
(
    id          INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title       VARCHAR(255)    NOT NULL,
    description VARCHAR(255),
    done        BOOLEAN         NOT NULL DEFAULT FALSE
);

CREATE TABLE lists_items
(
    id      INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    item_id INT             NOT NULL,
    list_id INT             NOT NULL,
    FOREIGN KEY (item_id) REFERENCES todo_items (id),
    FOREIGN KEY (list_id) REFERENCES todo_lists (id)
);