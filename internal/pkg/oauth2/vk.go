package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/renju24/backend/internal/pkg/config"
)

type vkResponse struct {
	Response []vkUserInfo `json:"response"`
}

type vkUserInfo struct {
	VkID     int64   `json:"id"`
	Username *string `json:"screen_name"`
}

func VKOauth(providers config.OauthConfig, code string, service Service, platform Platform) (*User, error) {
	oauthConfig, err := OauthConfig(providers, service, platform)
	if err != nil {
		return nil, err
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	var email *string
	if emailExtra := token.Extra("email"); emailExtra != nil {
		v, ok := emailExtra.(string)
		if ok {
			email = &v
		}
	}
	response, err := http.Get(providers.VK.API + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var vkResp vkResponse
	if err = json.Unmarshal(data, &vkResp); err != nil {
		return nil, err
	}
	if len(vkResp.Response) == 0 {
		return nil, errors.New("empty result")
	}
	vkUser := vkResp.Response[0]
	vkUserID := strconv.FormatInt(vkUser.VkID, 10)
	username := "vk" + vkUserID
	if vkUser.Username != nil {
		username = *vkUser.Username
	}
	return &User{
		ID:       vkUserID,
		Username: username,
		Email:    email,
	}, nil
}
