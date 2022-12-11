package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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
	GameID int64 `json:"game_id"`
}

func (apiServer *APIServer) CallForGame(c *websocket.Client, jsonData []byte) (*RPCCallForGameResponse, error) {
	var req RPCCallForGameRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, apierror.ErrorBadRequest
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
	ok, err := apiServer.db.IsPlaying(inviterID)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	// If user is already playing a game.
	if ok {
		return nil, apierror.ErrorAlreadyPlaying
	}
	opponent, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, apierror.ErrorUserNotFound
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	// Creating game in database with random black and white user and retrieve the game id.
	gameID, err := apiServer.db.CreateGame(randomBlackAndWhite(inviterID, opponent.ID))
	if err != nil {
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
	return &RPCCallForGameResponse{
		GameID: gameID,
	}, nil
}

var randUser = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomBlackAndWhite(user1, user2 int64) (blackUserID, whiteUserID int64) {
	switch randUser.Intn(2) {
	case 0:
		return user1, user2
	case 1:
		return user2, user1
	}
	return user1, user2
}
