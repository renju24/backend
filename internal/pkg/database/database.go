package database

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(dsn string) (*Database, error) {
	db, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		return nil, err
	}
	return &Database{
		pool: db,
	}, nil
}

func (db *Database) Close() error {
	if db.pool != nil {
		db.pool.Close()
	}
	return nil
}

func (db *Database) ReadConfig() (*config.Config, error) {
	query := `SELECT id, config_json FROM config ORDER BY id DESC LIMIT 1;`
	var (
		configVersion int
		configJSON    []byte
		config        config.Config
	)
	if err := db.pool.QueryRow(context.TODO(), query).Scan(&configVersion, &configJSON); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}
	config.Version = configVersion
	return &config, nil
}

func (db *Database) InsertUser(username, email, passwordBcrypt string) (userID int64, err error) {
	query := `INSERT INTO users (username, email, password_bcrypt, ranking) VALUES ($1, $2, $3, 400) RETURNING id;`
	if err = db.pool.QueryRow(context.TODO(), query, username, email, passwordBcrypt).Scan(&userID); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.ConstraintName == "unique_username" && pgxErr.Code == pgerrcode.UniqueViolation {
				return 0, apierror.ErrorUsernameIsTaken
			}
			if pgxErr.ConstraintName == "unique_email" && pgxErr.Code == pgerrcode.UniqueViolation {
				return 0, apierror.ErrorEmailIsTaken
			}
		}
		return 0, err
	}
	return userID, err
}

func (db *Database) GetLoginInfo(login string) (userID int64, passwordBcrypt string, err error) {
	query := "SELECT id, password_bcrypt FROM users "
	if strings.Contains(login, "@") {
		query += "WHERE email = $1"
	} else {
		query += "WHERE username = $1"
	}
	err = db.pool.QueryRow(context.TODO(), query, login).Scan(&userID, &passwordBcrypt)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", apierror.ErrorUserNotFound
	}
	return userID, passwordBcrypt, err
}
