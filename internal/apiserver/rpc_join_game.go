package apiserver

import (
	"encoding/json"

	"github.com/armantarkhanian/websocket"
)

type RPCJoinGameRequest struct{}

type RPCJoinGameResponse struct{}

func (app *APIServer) JoinGame(c *websocket.Client, jsonData []byte) (*RPCJoinGameResponse, error) {
	var req RPCJoinGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	return &RPCJoinGameResponse{}, nil
}
