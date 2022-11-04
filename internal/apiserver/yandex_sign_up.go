package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
)

type yandexUser struct {
	YandexID string `json:"id"`
	Username string `json:"login"`
	Email    string `json:"default_email"`
}

func yandexOauth(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthConfig, err := yandexOauthConfig(api, c.Param("platform"))
		if err != nil {
			log.Fatalln(1, err)
		}
		code := c.Request.FormValue("code")
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Fatalln(2, err)
		}
		req, err := http.NewRequest(http.MethodGet, api.config.Oauth2.Yandex.API, nil)
		if err != nil {
			log.Fatalln(10, err)
		}
		req.Header.Set("Authorization", "OAuth "+token.AccessToken)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(3, err)
		}
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(4, err)
		}
		var gUser yandexUser
		if err := json.Unmarshal(data, &gUser); err != nil {
			log.Fatalln(5, err)
		}
		user, err := api.db.CreateUserOauth(gUser.Username, gUser.Email, gUser.YandexID, config.OauthYandex)
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
		case api.config.Oauth2.Yandex.Callbacks.Web:
			c.SetCookie(
				api.config.Server.Token.Cookie.Name,
				jwtToken,
				api.config.Server.Token.Cookie.MaxAge,
				api.config.Server.Token.Cookie.Path,
				api.config.Server.Token.Cookie.Domain,
				api.config.Server.Token.Cookie.Secure,
				api.config.Server.Token.Cookie.HttpOnly,
			)
			api.logger.Info().
				Str("url", api.config.Oauth2.DeepLinks.Web).
				Msg("trying redirect to url")
			c.Redirect(http.StatusMovedPermanently, api.config.Oauth2.DeepLinks.Web)
			return
		case api.config.Oauth2.Yandex.Callbacks.Android:
			deepLink := api.config.Oauth2.DeepLinks.Android + "?token=" + jwtToken
			api.logger.Info().
				Str("url", deepLink).
				Msg("trying redirect to url")
			c.Redirect(http.StatusMovedPermanently, deepLink)
			return
		default:
			log.Fatalln("invalid URL fuck")
		}
	}
}
