package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/renju24/backend/internal/pkg/apierror"
	oauth "github.com/renju24/backend/internal/pkg/oauth2"
	"github.com/renju24/backend/model"
)

func (db *Database) createUserOauth(username string, email *string, oauthID string, service oauth.Service) (*model.User, error) {
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
					// if username is already taken, then select last username and increment it.
					var lastUsername string
					query = `SELECT username FROM users WHERE username LIKE $1 ORDER BY username DESC LIMIT 1;`
					if err = db.pool.QueryRow(ctx, query, username+"%").Scan(&lastUsername); err != nil {
						return nil, err
					}
					lastUsername = strings.TrimPrefix(lastUsername, username)
					lastUsername = strings.TrimPrefix(lastUsername, "-")
					lastUsername = strings.TrimSpace(lastUsername)
					var i int64
					if lastUsername != "" {
						i, err = strconv.ParseInt(lastUsername, 10, 64)
						if err != nil {
							return nil, err
						}
					}
					username = fmt.Sprintf("%s-%d", username, i+1)
					return db.createUserOauth(username, email, oauthID, service)
				case "unique_email":
					// TODO: what to do if email is already taken when using OAuth2 authorization.
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
