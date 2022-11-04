package apiserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/renju24/backend/internal/pkg/apierror"
	"golang.org/x/crypto/bcrypt"
)

type signinRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token"`
}

func signIn(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signinRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, &APIError{
				Error: apierror.ErrorBadRequest,
			})
			return
		}
		user, err := api.db.GetUserByLogin(req.Login)
		if err != nil {
			if errors.Is(err, apierror.ErrorUserNotFound) {
				c.JSON(http.StatusBadRequest, &APIError{
					Error: apierror.ErrorUserNotFound,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, &APIError{
				Error: apierror.ErrorInternal,
			})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordBcrypt), []byte(req.Password)) != nil {
			c.JSON(http.StatusBadRequest, &APIError{
				Error: apierror.ErrorInvalidCredentials,
			})
			return
		}

		jwtToken, err := api.jwt.Encode(jwt.Payload{
			Subject:        strconv.FormatInt(user.ID, 10),
			ExpirationTime: int64(api.config.Server.Token.Cookie.MaxAge),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, &APIError{
				Error: apierror.ErrorInternal,
			})
			return
		}

		resp := signinResponse{
			Token: jwtToken,
		}

		c.SetCookie(
			api.config.Server.Token.Cookie.Name,
			resp.Token,
			api.config.Server.Token.Cookie.MaxAge,
			api.config.Server.Token.Cookie.Path,
			api.config.Server.Token.Cookie.Domain,
			api.config.Server.Token.Cookie.Secure,
			api.config.Server.Token.Cookie.HttpOnly,
		)

		c.JSON(200, &resp)
	}
}
