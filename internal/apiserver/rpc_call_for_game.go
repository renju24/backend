package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
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
		return nil, apierror.ErrorInvalidBody
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return nil, apierror.ErrorUsernameIsRequired
	}
	inviterID, err := strconv.ParseInt(c.UserID(), 10, 64)
	if err != nil {
		return nil, apierror.ErrorUnauthorized
	}
	inviter, err := apiServer.db.GetUserByID(inviterID)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, apierror.ErrorUnauthorized
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	if inviter.Username == req.Username {
		return nil, apierror.ErrorCallingYourselfForGame
	}
	opponent, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, apierror.ErrorUserNotFound
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	_, err = apiServer.PublishEvent(fmt.Sprintf("user_%d", opponent.ID), &EventGameInvitation{
		Inviter:   inviter.Username,
		InvitedAt: time.Now(),
	})
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	// TODO: send push notifications.
	return &RPCCallForGameResponse{"ok"}, nil
}
