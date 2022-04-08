DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS tournament_types;

CREATE TABLE tournament_types (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL
);

CREATE TABLE users (
    id varchar(100) PRIMARY KEY,
    name varchar(100) NOT NULL,
    email varchar(100) NOT NULL,
    picture varchar(255) NOT NULL
);

CREATE TABLE tournaments (
    id serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    datetime TIMESTAMP NOT NULL,
    tournament_type_id INT NOT NULL,
    FOREIGN KEY (tournament_type_id) REFERENCES tournament_types (id)
);

CREATE TABLE participants (
    tournament_id INT NOT NULL,
    user_id varchar(100) NOT NULL,
    FOREIGN KEY (tournament_id) REFERENCES tournaments (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY(tournament_id, user_id)
);

INSERT INTO tournament_types (title) VALUES 
    ('Single elimination'),
    ('Double elimination'),
    ('Straight Round Robin'),
    ('Split round robin followed by single elimination');

INSERT INTO tournaments (title, description, datetime, tournament_type_id) VALUES
    ('Mario Kart Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-22 18:30', 4),
    ('Super Smash Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-09 19:00', 4);

INSERT INTO users values
    ('1', 'John McKelly', 'temp@example.com', 'http://www.url.com'),
    ('2', 'Sara Jonsson', 'temp@example.com', 'http://www.url.com'),
    ('3', 'Peter Smith', 'temp@example.com', 'http://www.url.com'),
    ('4', 'Lisa Clarksson', 'temp@example.com', 'http://www.url.com'),
    ('5', 'John Persson', 'temp@example.com', 'http://www.url.com'),
    ('6', 'Fredrik Lindberg', 'temp@example.com', 'http://www.url.com'),
    ('7', 'Molly Sandén', 'temp@example.com', 'http://www.url.com'),
    ('8', 'Kristian Luuk', 'temp@example.com', 'http://www.url.com'),
    ('9', 'Babben Larsson', 'temp@example.com', 'http://www.url.com'),
    ('10', 'David Sundin', 'temp@example.com', 'http://www.url.com');

INSERT INTO participants values 
    (2, 1),
    (2, 2),
    (2, 3),
    (2, 4),
    (2, 5),
    (2, 6),
    (2, 7),
    (2, 8),
    (2, 9),
    (2, 10);