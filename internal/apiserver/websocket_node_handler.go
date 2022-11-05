package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/renju24/backend/internal/pkg/apierror"
)

func (*APIServer) OnSurvey(*centrifuge.Node, centrifuge.SurveyEvent) centrifuge.SurveyReply {
	return centrifuge.SurveyReply{}
}

func (*APIServer) OnNotification(*centrifuge.Node, centrifuge.NotificationEvent) {}

func (apiServer *APIServer) OnConnecting(_ *centrifuge.Node, e centrifuge.ConnectEvent) (websocket.Session, centrifuge.ConnectReply, error) {
	if e.Token == "" {
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}

	var payload jwt.Payload
	if err := apiServer.jwt.Decode(e.Token, &payload); err != nil || payload.Subject == "" {
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}

	userID, err := strconv.ParseInt(payload.Subject, 10, 64)
	if err != nil {
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}

	user, err := apiServer.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
	}

	b, err := json.Marshal(&RPCGetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Ranking:  user.Ranking,
	})
	if err != nil {
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
	}

	websocketSession := WebsocketSession{
		UserID: userID,
	}

	apiServer.logger.Info().Int64("user_id", userID).Msg("user connected successfully")

	userIDString := fmt.Sprintf("%d", user.ID)

	if err := apiServer.centrifugeNode.Subscribe(userIDString, "user_"+userIDString); err != nil {
		return nil, centrifuge.ConnectReply{}, centrifuge.ErrorInternal
	}

	return &websocketSession, centrifuge.ConnectReply{Data: b}, nil
}

type WebsocketSession struct {
	UserID int64 `json:"user_id"`
}

func (ws *WebsocketSession) Authorized() bool {
	return ws.UserID != 0
}

func (ws *WebsocketSession) Credentials() *centrifuge.Credentials {
	return &centrifuge.Credentials{
		UserID: strconv.FormatInt(ws.UserID, 10),
	}
}
