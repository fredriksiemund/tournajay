DROP TABLE IF EXISTS tournament_types;

CREATE TABLE tournament_types (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL
);

DROP TABLE IF EXISTS tournaments;

CREATE TABLE tournaments (
    id serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    date TIMESTAMP NOT NULL,
    tournament_type INT NOT NULL,
    FOREIGN KEY (tournament_type) REFERENCES tournament_types (id)
);

INSERT INTO tournament_types (title) VALUES 
    ('Single elimination'),
    ('Double elimination'),
    ('Round robin'),
    ('Round robin and single elimination');

INSERT INTO tournaments (title, date, tournament_type) VALUES
    ('Mario Kart Tournament', '2022-05-22 18:30', 1),
    ('Super Smash Tournament', '2022-05-09 19:00', 3);
