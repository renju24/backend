package apiserver

import (
	"encoding/json"

	"github.com/armantarkhanian/websocket"
)

type RPCGetUserByIDRequest struct {
	UserID int64 `json:"user_id"`
}

type RPCGetUserByUsernameRequest struct {
	Username string `json:"username"`
}

type RPCGetUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Ranking  int    `json:"ranking"`
}

func (apiServer *APIServer) GetUserByID(c *websocket.Client, jsonData []byte) (*RPCGetUserResponse, error) {
	var req RPCGetUserByIDRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	user, err := apiServer.db.GetUserByID(req.UserID)
	if err != nil {
		return nil, err
	}
	return &RPCGetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Ranking:  user.Ranking,
	}, nil
}

func (apiServer *APIServer) GetUserByUsername(c *websocket.Client, jsonData []byte) (*RPCGetUserResponse, error) {
	var req RPCGetUserByUsernameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	user, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		return nil, err
	}
	return &RPCGetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Ranking:  user.Ranking,
	}, nil
}
