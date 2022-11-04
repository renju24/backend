package apiserver

import (
	"encoding/json"
	"strings"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/model"
)

type RPCGameHistoryRequest struct {
	Username string `json:"username"`
}

type RPCGameHistoryResponse struct {
	Games []model.GameHistoryItem `json:"games"`
}

func (app *APIServer) GameHistory(c *websocket.Client, jsonData []byte) (*RPCGameHistoryResponse, error) {
	var req RPCGameHistoryRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, apierror.ErrorInvalidBody
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return nil, apierror.ErrorUsernameIsRequired
	}
	games, err := app.db.GameHistory(req.Username)
	if err != nil {
		app.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	return &RPCGameHistoryResponse{
		Games: games,
	}, nil
}
