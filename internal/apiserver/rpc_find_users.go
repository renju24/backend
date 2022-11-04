package apiserver

import (
	"encoding/json"
	"strings"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCFindUserRequest struct {
	Username string `json:"username"`
}

type RPCFindUserResponse struct {
	Users []findUser `json:"users"`
}

type findUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Ranking  int    `json:"ranking"`
}

func (app *APIServer) FindUsers(_ *websocket.Client, jsonData []byte) (*RPCFindUserResponse, error) {
	var req RPCFindUserRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, apierror.ErrorBadRequest
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return nil, apierror.ErrorUsernameIsRequired
	}
	users, err := app.db.FindUsers(req.Username)
	if err != nil {
		app.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	var response RPCFindUserResponse
	for _, user := range users {
		response.Users = append(response.Users, findUser{
			ID:       user.ID,
			Username: user.Username,
			Ranking:  user.Ranking,
		})
	}
	return &response, nil
}
