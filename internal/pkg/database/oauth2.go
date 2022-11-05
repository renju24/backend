package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/renju24/backend/internal/pkg/apierror"
	oauth "github.com/renju24/backend/internal/pkg/oauth2"
	"github.com/renju24/backend/model"
)

func (db *Database) createUserOauth(username string, email *string, oauthID string, service oauth.Service, i int) (*model.User, error) {
	query := `INSERT INTO users (username, email, %s) VALUES ($1, $2, $3) RETURNING id, ranking;`
	switch service {
	case oauth.Google:
		query = fmt.Sprintf(query, "google_id")
	case oauth.Yandex:
		query = fmt.Sprintf(query, "yandex_id")
	case oauth.VK:
		query = fmt.Sprintf(query, "vk_id")
	default:
		return nil, oauth.ErrUnknownService
	}
	if i > 0 {
		username = fmt.Sprintf("%s-%d", username, i)
	}
	user := model.User{
		Username: username,
		Email:    email,
		GoogleID: &oauthID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	if err := db.pool.QueryRow(ctx, query, username, email, oauthID).Scan(&user.ID, &user.Ranking); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == pgerrcode.UniqueViolation {
				switch pgxErr.ConstraintName {
				case "unique_google_id":
					return db.getUserByOauthUserID(oauthID, oauth.Google)
				case "unique_yandex_id":
					return db.getUserByOauthUserID(oauthID, oauth.Yandex)
				case "unique_vk_id":
					return db.getUserByOauthUserID(oauthID, oauth.VK)
				case "unique_username":
					// if username is already taken then increment it.
					return db.createUserOauth(username, email, oauthID, service, i+1)
				case "unique_email":
					return nil, apierror.ErrorEmailIsTaken
				}
			}
		}
		return nil, err
	}
	return &user, nil
}

func (db *Database) getUserByOauthUserID(oauthID string, service oauth.Service) (*model.User, error) {
	query := `SELECT id, username, email, ranking, password_bcrypt FROM users WHERE %s = $1`
	switch service {
	case oauth.Google:
		query = fmt.Sprintf(query, "google_id")
	case oauth.Yandex:
		query = fmt.Sprintf(query, "yandex_id")
	case oauth.VK:
		query = fmt.Sprintf(query, "vk_id")
	default:
		return nil, oauth.ErrUnknownService
	}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultQueryTimeout)
	defer cancel()
	var user model.User
	err := db.pool.QueryRow(ctx, query, oauthID).Scan(
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
