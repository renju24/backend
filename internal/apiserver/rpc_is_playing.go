package apiserver

import (
	"encoding/json"
	"strconv"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCIsPlayingRequest struct{}

type RPCIsPlayingResponse struct {
	Game *playingGame `json:"game"`
}

type playingGame struct {
	GameID   int64  `json:"game_id"`
	Color    string `json:"color"`
	Opponent string `json:"opponent"`
}

func (apiServer *APIServer) IsPlaying(c *websocket.Client, jsonData []byte) (*RPCIsPlayingResponse, error) {
	var req RPCIsPlayingRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorBadRequest
	}
	userID, err := strconv.ParseInt(c.UserID(), 10, 64)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorUnauthorized
	}
	game, err := apiServer.db.GetPlayingGame(userID)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	if game == nil {
		return &RPCIsPlayingResponse{}, nil
	}
	user, err := apiServer.db.GetUserByID(userID)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	var (
		color    = "black"
		opponent = game.BlackUsername
	)
	if game.WhiteUsername == user.Username {
		color = "white"
	}
	if game.BlackUsername == user.Username {
		opponent = game.WhiteUsername
	}
	return &RPCIsPlayingResponse{
		Game: &playingGame{
			GameID:   game.ID,
			Color:    color,
			Opponent: opponent,
		},
	}, nil
}
