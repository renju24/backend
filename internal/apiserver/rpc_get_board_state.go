package apiserver

import (
	"encoding/json"
	"strconv"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCBoardStateRequest struct {
	GameID int64 `json:"game_id"`
}

type RPCBoardStateResponse struct {
	Moves []EventMove `json:"moves"`
}

func (app *APIServer) BoardState(c *websocket.Client, jsonData []byte) (*RPCBoardStateResponse, error) {
	var req RPCBoardStateRequest
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
	var response RPCBoardStateResponse
	for _, move := range moves {
		response.Moves = append(response.Moves, EventMove{
			UserID:      move.UserID,
			XCoordinate: move.XCoordinate,
			YCoordinate: move.YCoordinate,
		})
	}
	return &response, nil
}
