package database

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
	"github.com/renju24/backend/model"
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

func (db *Database) CreateUser(username, email, passwordBcrypt string) (*model.User, error) {
	user := model.User{
		Username:       username,
		Email:          email,
		PasswordBcrypt: passwordBcrypt,
	}
	query := `INSERT INTO users (username, email, password_bcrypt) VALUES ($1, $2, $3) RETURNING id, ranking;`
	if err := db.pool.QueryRow(context.TODO(), query, username, email, passwordBcrypt).Scan(
		&user.ID,
		&user.Ranking,
	); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.ConstraintName == "unique_username" && pgxErr.Code == pgerrcode.UniqueViolation {
				return nil, apierror.ErrorUsernameIsTaken
			}
			if pgxErr.ConstraintName == "unique_email" && pgxErr.Code == pgerrcode.UniqueViolation {
				return nil, apierror.ErrorEmailIsTaken
			}
		}
		return nil, err
	}
	return &user, nil
}

func (db *Database) GetUserByLogin(login string) (*model.User, error) {
	query := "SELECT id, username, email, ranking, password_bcrypt FROM users "
	if strings.Contains(login, "@") {
		query += "WHERE email = $1"
	} else {
		query += "WHERE username = $1"
	}
	var user model.User
	err := db.pool.QueryRow(context.TODO(), query, login).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Ranking,
		&user.PasswordBcrypt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apierror.ErrorUserNotFound
	}
	return &user, err
}

func (db *Database) GetUserByID(userID int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, ranking, password_bcrypt FROM users WHERE id = $1`
	err := db.pool.QueryRow(context.TODO(), query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Ranking,
		&user.PasswordBcrypt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apierror.ErrorUserNotFound
	}
	return &user, err
}

func (db *Database) CreateGame(blackUserID, whiteUserID int64) (*model.Game, error) {
	now := time.Now()
	game := model.Game{
		BlackUserID: blackUserID,
		WhiteUserID: whiteUserID,
		StartedAt:   now,
	}
	query := `INSERT INTO games (black_user_id, white_user_id, started_at) VALUES ($1, $2, $3) RETURNING id;`
	if err := db.pool.QueryRow(context.TODO(), query, blackUserID, whiteUserID, now).Scan(&game.ID); err != nil {
		return nil, err
	}
	return &game, nil
}

func (db *Database) IsGameMember(userID, gameID int64) (bool, error) {
	var ok bool
	query := `SELECT TRUE FROM games WHERE id = $1 AND (black_user_id = $2 OR white_user_id = $3)`
	if err := db.pool.QueryRow(context.TODO(), query, gameID, userID, userID).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (db *Database) FindUsers(username string) ([]*model.User, error) {
	username = strings.Trim(username, "%")
	username = "%" + username + "%"
	query := `SELECT id, username, email, ranking, password_bcrypt FROM users WHERE username ILIKE $1`
	rows, err := db.pool.Query(context.TODO(), query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Ranking, &user.PasswordBcrypt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, err
}
