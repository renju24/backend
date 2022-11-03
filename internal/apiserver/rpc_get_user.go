package apiserver

import (
	"encoding/json"
	"strconv"

	"github.com/armantarkhanian/websocket"
)

type RPCGetUserRequest struct {
	Username string `json:"username"`
}

type RPCGetUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Ranking  int    `json:"ranking"`
}

func (apiServer *APIServer) GetUser(c *websocket.Client, jsonData []byte) (*RPCGetUserResponse, error) {
	var req RPCGetUserRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	user, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		return nil, err
	}
	resp := RPCGetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Ranking:  user.Ranking,
	}
	if strconv.FormatInt(user.ID, 10) != c.UserID() {
		resp.Email = ""
	}
	return &resp, nil
}
