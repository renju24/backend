package apiserver

import (
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
)

func (*APIServer) OnAlive(*websocket.Client) {}

func (*APIServer) OnDisconect(*websocket.Client, centrifuge.DisconnectEvent) {}

func (*APIServer) OnSubscribe(*websocket.Client, centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
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
