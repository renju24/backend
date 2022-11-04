package apiserver

import (
	"github.com/armantarkhanian/websocket"
	"github.com/renju24/backend/internal/pkg/apierror"
)

func (apiServer *APIServer) Top10(_ *websocket.Client, _ []byte) (*RPCFindUserResponse, error) {
	users, err := apiServer.db.Top10()
	if err != nil {
		apiServer.logger.Error().Err(err).Send()
		return nil, apierror.ErrorInternal
	}
	var response RPCFindUserResponse
	for _, user := range users {
		response.Users = append(response.Users, findUser{
			ID:       user.ID,
			Username: user.Username,
			Ranking:  user.Ranking,
		})
	}
	return &response, nil
}
