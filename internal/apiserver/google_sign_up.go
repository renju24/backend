package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/armantarkhanian/jwt"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
)

type googleUser struct {
	GoogleID string `json:"id"`
	Email    string `json:"email"`
}

func googleOauth(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthConfig, err := googleOauthConfig(api, c.Param("platform"))
		if err != nil {
			log.Fatalln(err)
		}
		code := c.Request.FormValue("code")
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Fatalln(err)
		}
		response, err := http.Get(api.config.Oauth2.Google.API + token.AccessToken)
		if err != nil {
			log.Fatalln(err)
		}
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var gUser googleUser
		if err := json.Unmarshal(data, &gUser); err != nil {
			log.Fatalln(err)
		}
		username := strings.TrimSuffix(gUser.Email, "@gmail.com")
		user, err := api.db.CreateUserOauth(username, gUser.Email, gUser.GoogleID, config.OauthGoogle)
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
		case api.config.Oauth2.Google.Callbacks.Web:
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
		case api.config.Oauth2.Google.Callbacks.Android:
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
