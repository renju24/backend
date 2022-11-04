package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

const DefaultQueryTimeout = 5 * time.Second

type Database struct {
	pool *pgxpool.Pool
}

func New(dsn string) (*Database, error) {
	db, err := pgxpool.New(context.Background(), dsn)
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
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query).Scan(&configVersion, &configJSON); err != nil {
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
		PasswordBcrypt: &passwordBcrypt,
	}
	query := `INSERT INTO users (username, email, password_bcrypt) VALUES ($1, $2, $3) RETURNING id, ranking;`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, username, email, passwordBcrypt).Scan(
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

func (db *Database) CreateUserOauth(username, email, oauthID string, oauthSerivce config.OauthService) (*model.User, error) {
	switch oauthSerivce {
	case config.Google:
		return db.createGoogleUser(username, email, oauthID, 0)
	case config.Yandex:
		return db.createYandexUser(username, email, oauthID, 0)
	}
	return nil, errors.New("invalid oauth service")
}

func (db *Database) createGoogleUser(username, email, googleID string, i int) (*model.User, error) {
	if i > 0 {
		username = fmt.Sprintf("%s-%d", username, i)
	}
	user := model.User{
		Username: username,
		Email:    email,
		GoogleID: &googleID,
	}
	query := `INSERT INTO users (username, email, google_id) VALUES ($1, $2, $3) RETURNING id, ranking;`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, username, email, googleID).Scan(&user.ID, &user.Ranking); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.ConstraintName == "unique_google_id" && pgxErr.Code == pgerrcode.UniqueViolation {
				// if user with this google id already exists, then select and return him.
				return db.GetUserByGoogleID(googleID)
			}
			if pgxErr.ConstraintName == "unique_username" && pgxErr.Code == pgerrcode.UniqueViolation {
				// if username is already taken then increment it.
				return db.createGoogleUser(username, email, googleID, i+1)
			}
			if pgxErr.ConstraintName == "unique_email" && pgxErr.Code == pgerrcode.UniqueViolation {
				return nil, apierror.ErrorEmailIsTaken
			}
		}
		return nil, err
	}
	return &user, nil
}

func (db *Database) createYandexUser(username, email, yandexID string, i int) (*model.User, error) {
	if i > 0 {
		username = fmt.Sprintf("%s-%d", username, i)
	}
	user := model.User{
		Username: username,
		Email:    email,
		YandexID: &yandexID,
	}
	query := `INSERT INTO users (username, email, yandex_id) VALUES ($1, $2, $3) RETURNING id, ranking;`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, username, email, yandexID).Scan(&user.ID, &user.Ranking); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.ConstraintName == "unique_yandex_id" && pgxErr.Code == pgerrcode.UniqueViolation {
				// if user with this yandex id already exists, then select and return him.
				return db.GetUserByYandexID(yandexID)
			}
			if pgxErr.ConstraintName == "unique_username" && pgxErr.Code == pgerrcode.UniqueViolation {
				// if username is already taken then increment it.
				return db.createYandexUser(username, email, yandexID, i+1)
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
	query := "SELECT id, username, email, google_id, yandex_id, ranking, password_bcrypt FROM users "
	if strings.Contains(login, "@") {
		query += "WHERE email = $1"
	} else {
		query += "WHERE username = $1"
	}
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	err := db.pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.GoogleID,
		&user.YandexID,
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
	query := `SELECT id, username, email, google_id, yandex_id, ranking, password_bcrypt FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	err := db.pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.GoogleID,
		&user.YandexID,
		&user.Ranking,
		&user.PasswordBcrypt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apierror.ErrorUserNotFound
	}
	return &user, err
}

func (db *Database) GetUserByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, google_id, yandex_id, ranking, password_bcrypt FROM users WHERE google_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	err := db.pool.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.GoogleID,
		&user.YandexID,
		&user.Ranking,
		&user.PasswordBcrypt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apierror.ErrorUserNotFound
	}
	return &user, err
}

func (db *Database) GetUserByYandexID(yandexID string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, google_id, yandex_id, ranking, password_bcrypt FROM users WHERE yandex_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	err := db.pool.QueryRow(ctx, query, yandexID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.GoogleID,
		&user.YandexID,
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
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, blackUserID, whiteUserID, now).Scan(&game.ID); err != nil {
		return nil, err
	}
	return &game, nil
}

func (db *Database) IsGameMember(userID, gameID int64) (bool, error) {
	var ok bool
	query := `SELECT TRUE FROM games WHERE id = $1 AND (black_user_id = $2 OR white_user_id = $3)`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, gameID, userID, userID).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (db *Database) FindUsers(username string) ([]*model.User, error) {
	username = strings.Trim(username, "%")
	username = "%" + username + "%"
	query := `SELECT id, username, email, ranking, password_bcrypt FROM users WHERE username ILIKE $1`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	rows, err := db.pool.Query(ctx, query, username)
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

func (db *Database) GameHistory(username string) ([]model.GameHistoryItem, error) {
	query := `
		SELECT
			g.id,
			black.username as black_username,
			white.username as white_username,
			winner.username as winner
		FROM
			games g
			INNER JOIN users black ON g.black_user_id = black.id
			INNER JOIN users white ON g.white_user_id = white.id
			LEFT  JOIN users winner ON g.winner_id = winner.id
		WHERE
			g.finished_at IS NOT NULL
			AND (black.username = $1 OR white.username = $1);`
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	rows, err := db.pool.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var games []model.GameHistoryItem
	for rows.Next() {
		var game model.GameHistoryItem
		if err = rows.Scan(&game.ID, &game.BlackUsername, &game.WhiteUsername, &game.WinnerUsername); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func (db *Database) Top10() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	rows, err := db.pool.Query(ctx, `SELECT id, username, ranking FROM users ORDER BY ranking DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Username, &user.Ranking); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, err
}
