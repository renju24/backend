package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/renju24/backend/internal/pkg/config"
)

type googleUserInfo struct {
	GoogleID string `json:"id"`
	Email    string `json:"email"`
}

func GoogleOauth(providers config.OauthConfig, code string, service Service, platform Platform) (*User, error) {
	oauthConfig, err := OauthConfig(providers, service, platform)
	if err != nil {
		return nil, err
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(providers.Google.API + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var googleUser googleUserInfo
	if err = json.Unmarshal(data, &googleUser); err != nil {
		return nil, err
	}
	return &User{
		ID:       googleUser.GoogleID,
		Username: strings.TrimSuffix(googleUser.Email, "@gmail.com"),
		Email:    &googleUser.Email,
	}, nil
}
