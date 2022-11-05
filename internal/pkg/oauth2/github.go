package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/renju24/backend/internal/pkg/config"
)

type githubUserInfo struct {
	GithubID int64   `json:"id"`
	Username string  `json:"login"`
	Email    *string `json:"email"`
}

func GithubOauth(providers config.OauthConfig, code string, service Service, platform Platform) (*User, error) {
	oauthConfig, err := OauthConfig(providers, service, platform)
	if err != nil {
		return nil, err
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, providers.Github.API, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var githubUser githubUserInfo
	if err = json.Unmarshal(data, &githubUser); err != nil {
		return nil, err
	}
	githubID := strconv.FormatInt(githubUser.GithubID, 10)
	return &User{
		ID:       githubID,
		Username: githubUser.Username,
		Email:    githubUser.Email,
	}, nil
}
