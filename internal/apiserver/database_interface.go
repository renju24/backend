package apiserver

import "github.com/renju24/backend/apimodel"

type Database interface {
	// Insert user into database.
	InsertUser(username, email, passwordBcrypt string) (userID int64, err error)
	// Get userID and password by login.
	GetLoginInfo(login string) (userID int64, passwordBcrypt string, err error)
	// Get user.
	GetUser(userID int64) (apimodel.User, error)
	// Close
	Close() error
}
