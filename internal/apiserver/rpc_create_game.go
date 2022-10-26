package apiserver

import (
	"encoding/json"

	"github.com/armantarkhanian/websocket"
)

type RPCCreateGameRequest struct{}

type RPCCreateGameResponse struct{}

func (app *APIServer) CreateGame(c *websocket.Client, jsonData []byte) (*RPCCreateGameResponse, error) {
	var req RPCCreateGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	return &RPCCreateGameResponse{}, nil
}
