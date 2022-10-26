package apiserver

import (
	"strconv"
	"strings"

	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
)

func (*APIServer) OnAlive(*websocket.Client) {}

func (*APIServer) OnDisconect(*websocket.Client, centrifuge.DisconnectEvent) {}

func (app *APIServer) OnSubscribe(c *websocket.Client, e centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	if strings.HasPrefix(e.Channel, "user_") {
		if e.Channel != "user_"+c.UserID() {
			return centrifuge.SubscribeReply{}, centrifuge.ErrorPermissionDenied
		}
	}

	if strings.HasPrefix(e.Channel, "game_") {
		gameID, err := strconv.ParseInt(e.Channel[5:], 10, 64)
		if err != nil {
			return centrifuge.SubscribeReply{}, centrifuge.ErrorBadRequest
		}
		userID, err := strconv.ParseInt(c.UserID(), 10, 64)
		if err != nil {
			return centrifuge.SubscribeReply{}, centrifuge.ErrorBadRequest
		}

		// Check if the user is a game member.
		ok, err := app.db.IsGameMember(gameID, userID)
		if err != nil {
			return centrifuge.SubscribeReply{}, centrifuge.ErrorInternal
		}
		if !ok {
			return centrifuge.SubscribeReply{}, centrifuge.ErrorPermissionDenied
		}
	}

	return centrifuge.SubscribeReply{}, nil
}

func (*APIServer) OnUnsubscribe(*websocket.Client, centrifuge.UnsubscribeEvent) {}

func (*APIServer) OnPublish(*websocket.Client, centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	return centrifuge.PublishReply{}, nil
}

func (*APIServer) OnRefresh(*websocket.Client, centrifuge.RefreshEvent) (centrifuge.RefreshReply, error) {
	return centrifuge.RefreshReply{}, nil
}

func (*APIServer) OnSubRefresh(*websocket.Client, centrifuge.SubRefreshEvent) (centrifuge.SubRefreshReply, error) {
	return centrifuge.SubRefreshReply{}, nil
}

func (*APIServer) OnMessage(*websocket.Client, centrifuge.MessageEvent) {}

func (*APIServer) OnPresence(*websocket.Client, centrifuge.PresenceEvent) (centrifuge.PresenceReply, error) {
	return centrifuge.PresenceReply{}, nil
}

func (*APIServer) OnPresenceStats(*websocket.Client, centrifuge.PresenceStatsEvent) (centrifuge.PresenceStatsReply, error) {
	return centrifuge.PresenceStatsReply{}, nil
}

func (*APIServer) OnRPC(*websocket.Client, centrifuge.RPCEvent) (centrifuge.RPCReply, error) {
	return centrifuge.RPCReply{}, nil
}

func (*APIServer) OnHistory(*websocket.Client, centrifuge.HistoryEvent) (centrifuge.HistoryReply, error) {
	return centrifuge.HistoryReply{}, nil
}
