package database

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) (*Database, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func (db *Database) Close() error {
	if db.db != nil {
		return db.db.Close()
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
	if err := db.db.QueryRow(query).Scan(&configVersion, &configJSON); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}
	config.Version = configVersion
	return &config, nil
}

func (db *Database) InsertUser(username string, email string, passwordBcrypt string) (userID int64, err error) {
	query := `INSERT INTO users (username, email, password_bcrypt, ranking) VALUES ($1, $2, $3, 400) RETURNING id;`
	err = db.db.QueryRow(query, username, email, passwordBcrypt).Scan(&userID)
	if err != nil {
		if pgxErr, ok := err.(*pgconn.PgError); ok {
			if pgxErr.ColumnName == "username" && pgxErr.Code == pgerrcode.UniqueViolation {
				return 0, apierror.ErrorUsernameIsTaken
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
	err = db.db.QueryRow(query, login).Scan(&userID, &passwordBcrypt)
	return userID, passwordBcrypt, err
}
