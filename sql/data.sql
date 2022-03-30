DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS tournament_types;

CREATE TABLE tournament_types (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL
);

CREATE TABLE tournaments (
    id serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    datetime TIMESTAMP NOT NULL,
    tournament_type_id INT NOT NULL,
    FOREIGN KEY (tournament_type_id) REFERENCES tournament_types (id)
);

INSERT INTO tournament_types (title) VALUES 
    ('Single elimination'),
    ('Double elimination'),
    ('Straight Round Robin'),
    ('Split round robin followed by single elimination');

INSERT INTO tournaments (title, description, datetime, tournament_type_id) VALUES
    ('Mario Kart Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-22 18:30', 1),
    ('Super Smash Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-09 19:00', 3);
