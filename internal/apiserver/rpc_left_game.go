package apiserver

import (
	"encoding/json"

	"github.com/armantarkhanian/websocket"
)

type RPCLeftGameRequest struct{}

type RPCLeftGameResponse struct{}

func (app *APIServer) LeftGame(c *websocket.Client, jsonData []byte) (*RPCLeftGameResponse, error) {
	var req RPCLeftGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	return &RPCLeftGameResponse{}, nil
}
