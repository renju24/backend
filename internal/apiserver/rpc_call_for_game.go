package apiserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
)

type RPCCallForGameRequest struct {
	Username string `json:"username"`
}

type RPCCallForGameResponse struct {
	Status string `json:"status"`
}

func (apiServer *APIServer) CallForGame(c *websocket.Client, jsonData []byte) (*RPCCallForGameResponse, error) {
	var req RPCCallForGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, err
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return nil, centrifuge.ErrorBadRequest
	}
	inviterID, err := strconv.ParseInt(c.UserID(), 10, 64)
	if err != nil {
		return nil, centrifuge.ErrorInternal
	}
	inviter, err := apiServer.db.GetUserByID(inviterID)
	if err != nil {
		return nil, centrifuge.ErrorInternal
	}
	if inviter.Username == req.Username {
		return &RPCCallForGameResponse{"ok"}, nil
	}
	opponent, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		return nil, centrifuge.ErrorInternal
	}
	_, err = apiServer.PublishEvent(fmt.Sprintf("user_%d", opponent.ID), &EventGameInvitation{
		Inviter:   inviter.Username,
		InvitedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	// TODO: send push notifications.
	return &RPCCallForGameResponse{"ok"}, nil
}
