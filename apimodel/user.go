package apimodel

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Ranking  int    `json:"ranking"`
}
