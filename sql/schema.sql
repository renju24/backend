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
