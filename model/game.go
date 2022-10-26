package model

import "time"

type Game struct {
	ID          int64      `json:"id"`
	BlackUserID int64      `json:"black_user_id"`
	WhiteUserID int64      `json:"white_user_id"`
	Winner      *int64     `json:"winner_id"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}
