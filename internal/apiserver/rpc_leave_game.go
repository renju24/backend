package apiserver

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCLeaveGameRequest struct {
	GameID int64 `json:"game_id"`
}

type RPCLeaveGameResponse struct{}

func (app *APIServer) LeaveGame(c *websocket.Client, jsonData []byte) (*RPCLeaveGameResponse, error) {
	var req RPCLeaveGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, apierror.ErrorBadRequest
	}
	userID, err := strconv.ParseInt(c.UserID(), 10, 64)
	if err != nil {
		return nil, apierror.ErrorUnauthorized
	}
	isGameMember, err := app.db.IsGameMember(userID, req.GameID)
	if err != nil {
		app.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	if !isGameMember {
		return nil, apierror.ErrorPermissionDenied
	}
	moves, err := app.db.GetGameMovesByID(req.GameID)
	if err != nil {
		app.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	event := &EventUserLeftGame{
		WhoLeftGameID: userID,
	}
	if len(moves) >= 2 {
		// If there is more >= 2 moves, then user lose a game.
		game, err := app.db.GetGameByID(req.GameID)
		if err != nil {
			app.logger.Error().Err(err).Send()
			return nil, apierror.ErrorInternal
		}
		winnerID := game.BlackUserID
		if userID == game.BlackUserID {
			winnerID = game.WhiteUserID
		}
		if err = app.db.FinishGameWithWinner(req.GameID, winnerID); err != nil {
			app.logger.Error().Err(err).Send()
		}
		event.WinnerID = &winnerID
	} else {
		//  Otherwise, just delete a game.
		if err = app.db.DeleteGame(req.GameID); err != nil {
			app.logger.Error().Err(err).Send()
		}
	}
	// Publish event.
	gameChannel := fmt.Sprintf("game_%d", req.GameID)
	if _, err = app.PublishEvent(gameChannel, event); err != nil {
		app.logger.Error().Err(err).Send()
	}
	return &RPCLeaveGameResponse{}, nil
}
