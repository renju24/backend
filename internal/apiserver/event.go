package apiserver

import (
	"encoding/json"
	"time"

	"github.com/centrifugal/centrifuge"
)

type Event interface {
	EventType() string
}

func (apiServer *APIServer) PublishEvent(channel string, event Event) (centrifuge.PublishResult, error) {
	if event == nil {
		return centrifuge.PublishResult{}, nil
	}
	msg, err := json.Marshal(map[string]any{
		"event_type": event.EventType(),
		"data":       event,
	})
	if err != nil {
		return centrifuge.PublishResult{}, err
	}
	res, err := apiServer.centrifugeNode.Publish(channel, msg)
	if err != nil {
		return centrifuge.PublishResult{}, err
	}
	return res, nil
}

type EventGameInvitation struct {
	Inviter   string    `json:"inviter"`
	InvitedAt time.Time `json:"invited_at"`
}

func (e *EventGameInvitation) EventType() string {
	return "game_invitation"
}

type EventDeclineGameInvitation struct{}

func (e *EventDeclineGameInvitation) EventType() string {
	return "decline_game_invitation"
}

type EventGameStarted struct{}

func (e *EventGameStarted) EventType() string {
	return "game_started"
}

type EventGameInvitationExpired struct{}

func (e *EventGameInvitationExpired) EventType() string {
	return "game_invitation_expired"
}
