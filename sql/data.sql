DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS tournament_types;

CREATE TABLE tournament_types (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL
);

CREATE TABLE tournaments (
    id serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    datetime TIMESTAMP NOT NULL,
    tournament_type INT NOT NULL,
    FOREIGN KEY (tournament_type) REFERENCES tournament_types (id)
);

INSERT INTO tournament_types (title) VALUES 
    ('Single elimination'),
    ('Double elimination'),
    ('Round robin'),
    ('Round robin followed by single elimination');

INSERT INTO tournaments (title, datetime, tournament_type) VALUES
    ('Mario Kart Tournament', '2022-05-22 18:30', 1),
    ('Super Smash Tournament', '2022-05-09 19:00', 3);
