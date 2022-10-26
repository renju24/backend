package apiserver

import (
	"github.com/renju24/backend/model"
)

type Database interface {
	// Create new user.
	CreateUser(username, email, passwordBcrypt string) (*model.User, error)
	// Get user by login.
	GetUserByLogin(login string) (*model.User, error)
	// Get user by ID.
	GetUserByID(userID int64) (*model.User, error)

	// Create new game.
	CreateGame(blackUserID, whiteUserID int64) (*model.Game, error)

	// Close
	Close() error
}
