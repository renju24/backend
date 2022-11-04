package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCCallForGameRequest struct {
	Username string `json:"username"`
}

type RPCCallForGameResponse struct {
	Status string `json:"status"`
}

func (apiServer *APIServer) CallForGame(c *websocket.Client, jsonData []byte) (*RPCCallForGameResponse, *apierror.Error, *centrifuge.Error) {
	var req RPCCallForGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, apierror.ErrorInvalidBody, centrifuge.ErrorBadRequest
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return nil, apierror.ErrorUsernameIsRequired, centrifuge.ErrorBadRequest
	}
	inviterID, err := strconv.ParseInt(c.UserID(), 10, 64)
	if err != nil {
		return nil, apierror.ErrorUnauthorized, centrifuge.ErrorUnauthorized
	}
	inviter, err := apiServer.db.GetUserByID(inviterID)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, apierror.ErrorUnauthorized, centrifuge.ErrorUnauthorized
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal, centrifuge.ErrorInternal
	}
	if inviter.Username == req.Username {
		return nil, apierror.ErrorCallingYourselfForGame, centrifuge.ErrorBadRequest
	}
	opponent, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, apierror.ErrorUserNotFound, centrifuge.ErrorBadRequest
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal, centrifuge.ErrorInternal
	}
	_, err = apiServer.PublishEvent(fmt.Sprintf("user_%d", opponent.ID), &EventGameInvitation{
		Inviter:   inviter.Username,
		InvitedAt: time.Now(),
	})
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal, centrifuge.ErrorInternal
	}
	// TODO: send push notifications.
	return &RPCCallForGameResponse{"ok"}, nil, nil
}
