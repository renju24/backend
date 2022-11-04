package apiserver

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/renju24/backend/internal/pkg/apierror"
)

type RPCGetUserRequest struct {
	Username string `json:"username"`
}

type RPCGetUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Ranking  int    `json:"ranking"`
}

func (apiServer *APIServer) GetUser(c *websocket.Client, jsonData []byte) (*RPCGetUserResponse, *centrifuge.Error) {
	var req RPCGetUserRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		return nil, &centrifuge.Error{
			Code:      uint32(apierror.ErrorInvalidBody.Code),
			Message:   apierror.ErrorInvalidBody.Message,
			Temporary: false,
		}
	}
	user, err := apiServer.db.GetUserByLogin(req.Username)
	if err != nil {
		if errors.Is(err, apierror.ErrorUserNotFound) {
			return nil, &centrifuge.Error{
				Code:      uint32(apierror.ErrorUserNotFound.Code),
				Message:   apierror.ErrorUserNotFound.Message,
				Temporary: false,
			}
		}
		apiServer.logger.Error().Err(err).Send()
		return nil, &centrifuge.Error{
			Code:      uint32(apierror.ErrorInternal.Code),
			Message:   apierror.ErrorInternal.Message,
			Temporary: false,
		}
	}
	resp := RPCGetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Ranking:  user.Ranking,
	}
	if strconv.FormatInt(user.ID, 10) != c.UserID() {
		resp.Email = ""
	}
	return &resp, nil
}
