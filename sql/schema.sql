CREATE TABLE config (
	id          SERIAL PRIMARY KEY,
	config_json JSONB  NOT NULL
);

CREATE TABLE users (
	id              SERIAL       PRIMARY KEY,
	username        VARCHAR(32)  NOT NULL,
	email           VARCHAR(84)  NOT NULL,
	password_bcrypt VARCHAR(128) NOT NULL,
	ranking         INT          NOT NULL DEFAULT 400
);
CREATE UNIQUE INDEX unique_username ON users (username);
CREATE UNIQUE INDEX unique_email ON users (email);

CREATE TABLE games (
	id              SERIAL       PRIMARY KEY,
	black_user_id   INT          NOT NULL REFERENCES users(id),
	white_user_id   INT          NOT NULL REFERENCES users(id),
	winner_id       INT          NULL     REFERENCES users(id),
	started_at      TIMESTAMP(0) NULL,
	finished_at     TIMESTAMP(0) NULL
);

CREATE TABLE moves (
	game_id         INT NOT NULL REFERENCES games(id),
	user_id         INT NOT NULL REFERENCES users(id),
	x_coordinate    INT NOT NULL,
	y_coordinate    INT NOT NULL
);
CREATE INDEX moves_game_id ON moves (game_id);
