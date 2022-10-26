package apiserver

import (
	"encoding/json"
)

type EventType string

const (
	UserJoinGame EventType = "user_join_game"
	UserLeftGame EventType = "user_left_game"
)

type Event struct {
	EventType EventType `json:"event_type"`
	Data      any       `json:"data"`
}

type (
	EventJoinGame struct {
		UserID int64 `json:"user_id"`
		GameID int64 `json:"game_id"`
	}
	EventLeftGame struct {
		UserID int64 `json:"user_id"`
		GameID int64 `json:"game_id"`
	}
)

func (apiServer *APIServer) PublishEvent(channel string, event Event) {
	switch event.EventType {
	case UserJoinGame:
		if _, ok := event.Data.(*EventJoinGame); !ok {
			return
		}
	case UserLeftGame:
		if _, ok := event.Data.(*EventLeftGame); !ok {
			return
		}
	}
	msg, err := json.Marshal(event)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return
	}
	_, err = apiServer.centrifugeNode.Publish(channel, msg)
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return
	}
}
