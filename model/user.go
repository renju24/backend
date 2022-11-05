package model

type User struct {
	ID             int64   `json:"id"`
	Username       string  `json:"username"`
	Email          *string `json:"email"`
	Ranking        int     `json:"ranking"`
	PasswordBcrypt *string `json:"-"`
	GoogleID       *string `json:"-"`
	YandexID       *string `json:"-"`
	GithubID       *string `json:"-"`
	VKID           *string `json:"-"`
}
