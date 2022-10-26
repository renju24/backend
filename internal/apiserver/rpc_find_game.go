package apiserver

import (
	"encoding/json"

	"github.com/armantarkhanian/websocket"
)

type RPCFindGameRequest struct{}

type RPCFindGameResponse struct{}

func (app *APIServer) FindGame(c *websocket.Client, jsonData []byte) (*RPCFindGameResponse, error) {
	var req RPCFindGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	return &RPCFindGameResponse{}, nil
}
