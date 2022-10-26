package apiserver

type Database interface {
	// Insert user into database.
	InsertUser(username, email, passwordBcrypt string) (userID int64, err error)
	// Get userID and password by login.
	GetLoginInfo(login string) (userID int64, passwordBcrypt string, err error)
}
