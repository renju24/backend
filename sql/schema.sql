CREATE TABLE config (
	id          SERIAL PRIMARY KEY,
	config_json JSONB  NOT NULL
);

CREATE TABLE users (
	id              SERIAL       PRIMARY KEY,
	username        VARCHAR(32)  NOT NULL,
	email           VARCHAR(84)  NULL,
	password_bcrypt VARCHAR(128) NULL,
	google_id       VARCHAR(64)  NULL,
	yandex_id       VARCHAR(64)  NULL,
	vk_id           VARCHAR(64)  NULL,
	ranking         INT          NOT NULL DEFAULT 400
);
CREATE UNIQUE INDEX unique_google_id ON users (google_id);
CREATE UNIQUE INDEX unique_yandex_id ON users (yandex_id);
CREATE UNIQUE INDEX unique_vk_id ON users (vk_id);
CREATE UNIQUE INDEX unique_username ON users (username);
CREATE UNIQUE INDEX unique_email ON users (email);

CREATE TABLE games (
	id              SERIAL       PRIMARY KEY,
	black_user_id   INT          NOT NULL REFERENCES users(id),
	white_user_id   INT          NOT NULL REFERENCES users(id),
	winner_id       INT          NULL     REFERENCES users(id),
	status          INT          NOT NULL,
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
