DROP TABLE IF EXISTS game_paths;
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS result_types;
DROP TABLE IF EXISTS contestants;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tournament_types;

CREATE TABLE tournament_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    picture VARCHAR(255) NOT NULL
);

CREATE TABLE tournaments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    datetime TIMESTAMP NOT NULL,
    tournament_type_id INT NOT NULL,
    creator_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (tournament_type_id) REFERENCES tournament_types (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
);

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE participants (
    tournament_id INT,
    user_id VARCHAR(255),
    team_id INT,
    FOREIGN KEY (tournament_id) REFERENCES tournaments (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (team_id) REFERENCES teams (id),
    PRIMARY KEY (tournament_id, user_id)
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL,
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
);

CREATE TABLE contestants (
    game_id INT,
    team_id INT,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    PRIMARY KEY (game_id, team_id)
);

CREATE TABLE result_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE results (
    game_id INT,
    team_id INT,
    result_type_id INT,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (result_type_id) REFERENCES result_types(id),
    PRIMARY KEY (game_id, team_id, result_type_id)
);

CREATE TABLE game_paths (
    from_game_id INT,
    to_game_id INT,
    result_type_id INT,
    FOREIGN KEY (from_game_id) REFERENCES games(id),
    FOREIGN KEY (to_game_id) REFERENCES games(id),
    FOREIGN KEY (result_type_id) REFERENCES result_types(id),
    PRIMARY KEY (from_game_id, to_game_id, result_type_id)
);

-- STATIC DATA --
INSERT INTO tournament_types (name) VALUES 
    ('Single elimination'),
    ('Double elimination'),
    ('Straight Round Robin'),
    ('Split round robin followed by single elimination');

INSERT INTO result_types (name) VALUES
    ('First'),
    ('Second');

-- TEST DATA --
INSERT INTO users values
    ('1', 'John McKelly', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/1.jpg'),
    ('2', 'Sara Jonsson', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/2.jpg'),
    ('3', 'Peter Smith', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/3.jpg'),
    ('4', 'Lisa Clarksson', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/4.jpg'),
    ('5', 'John Persson', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/5.jpg'),
    ('6', 'Fredrik Lindberg', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/6.jpg'),
    ('7', 'Molly Sandén', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/7.jpg'),
    ('8', 'Kristian Luuk', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/8.jpg'),
    ('9', 'Babben Larsson', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/9.jpg'),
    ('10', 'David Sundin', 'temp@example.com', 'https://randomuser.me/api/portraits/lego/1.jpg');

INSERT INTO tournaments (title, description, datetime, tournament_type_id, creator_id) VALUES
    ('Mario Kart Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-22 18:30', 4, '1'),
    ('Super Smash Tournament', 'With supporting text below as a natural lead-in to additional content!', '2022-05-09 19:00', 4, '2');

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