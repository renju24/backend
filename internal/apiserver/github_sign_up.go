package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
)

type githubUserInfo struct {
	GithubID int64   `json:"id"`
	Username string  `json:"login"`
	Email    *string `json:"email"`
}

func githubOauth(api *APIServer, c *gin.Context, service config.OauthService, platform config.Platform) {
	oauthConfig, err := oauthConfig(api, service, platform)
	if err != nil {
		api.logger.Err(err).Send()
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	code := c.Request.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusBadRequest, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	req, err := http.NewRequest(http.MethodGet, api.config.Oauth2.Github.API, nil)
	if err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusBadRequest, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusBadRequest, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusBadRequest, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	var githubUser githubUserInfo
	if err = json.Unmarshal(data, &githubUser); err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusBadRequest, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	githubID := strconv.FormatInt(githubUser.GithubID, 10)
	user, err := api.db.CreateUserOauth(githubUser.Username, githubUser.Email, githubID, config.Github)
	if err != nil {
		if errors.Is(err, apierror.ErrorEmailIsTaken) {
			c.JSON(http.StatusBadRequest, &apierror.Error{
				Error: apierror.ErrorEmailIsTaken,
			})
			return
		}
		api.logger.Err(err).Send()
		c.JSON(http.StatusInternalServerError, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	jwtToken, err := api.jwt.Encode(jwt.Payload{
		Subject:        strconv.FormatInt(user.ID, 10),
		ExpirationTime: int64(api.config.Server.Token.Cookie.MaxAge),
	})
	if err != nil {
		api.logger.Err(err).Send()
		c.JSON(http.StatusInternalServerError, &apierror.Error{
			Error: apierror.ErrorInternal,
		})
		return
	}
	switch oauthConfig.RedirectURL {
	case api.config.Oauth2.Github.Callbacks.Android:
		deepLink := api.config.Oauth2.DeepLinks.Android + "?token=" + jwtToken
		c.Redirect(http.StatusFound, deepLink)
		return
	default:
		c.SetCookie(
			api.config.Server.Token.Cookie.Name,
			jwtToken,
			api.config.Server.Token.Cookie.MaxAge,
			api.config.Server.Token.Cookie.Path,
			api.config.Server.Token.Cookie.Domain,
			api.config.Server.Token.Cookie.Secure,
			api.config.Server.Token.Cookie.HttpOnly,
		)
		c.Redirect(http.StatusFound, api.config.Oauth2.DeepLinks.Web)
		return
	}
}
