package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/renju24/backend/internal/pkg/config"
)

type yandexUserInfo struct {
	YandexID string `json:"id"`
	Username string `json:"login"`
	Email    string `json:"default_email"`
}

func YandexOauth(providers config.OauthConfig, code string, service Service, platform Platform) (*User, error) {
	oauthConfig, err := OauthConfig(providers, service, platform)
	if err != nil {
		return nil, err
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, providers.Yandex.API, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+token.AccessToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var yandexUser yandexUserInfo
	if err = json.Unmarshal(data, &yandexUser); err != nil {
		return nil, err
	}
	return &User{
		ID:       yandexUser.YandexID,
		Username: yandexUser.Username,
		Email:    &yandexUser.Email,
	}, nil
}
