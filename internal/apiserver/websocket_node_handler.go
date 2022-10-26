package apiserver

import (
	"encoding/json"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/renju24/backend/apimodel"
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
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}

	b, err := json.Marshal(&user)
	if err != nil {
		return nil, centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
	}

	websocketSession := WebsocketSession{
		User: apimodel.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Ranking:  user.Ranking,
		},
	}

	return &websocketSession, centrifuge.ConnectReply{Data: b}, nil
}

type WebsocketSession struct {
	User apimodel.User
}

func (ws *WebsocketSession) Authorized() bool {
	return ws.User.ID != 0
}

func (ws *WebsocketSession) Credentials() *centrifuge.Credentials {
	return &centrifuge.Credentials{
		UserID: strconv.FormatInt(ws.User.ID, 10),
	}
}
