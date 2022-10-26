package apiserver

import (
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/renju24/backend/internal/pkg/apierror"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

type signupResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

func signUp(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signupRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, &APIError{
				Error: apierror.ErrorInvalidBody,
			})
			return
		}
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, &APIError{
				Error: apierror.ErrorInvalidBody,
			})
			return
		}
		passwordBcrypt, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &APIError{
				Error: apierror.ErrorInternal,
			})
			return
		}
		_, err = api.db.InsertUser(
			req.Username,
			req.Email,
			string(passwordBcrypt),
		)
		if err != nil {
			// TODO: check error if username is already taken.
			c.JSON(http.StatusInternalServerError, &APIError{
				Error: apierror.ErrorInternal,
			})
			return
		}

		resp := signupResponse{
			Status: 1,
			Token:  "", // TODO: generate JWT-token.
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

func (req *signupRequest) Validate() error {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" {
		return apierror.ErrorUsernameIsRequired
	}
	if req.Email == "" {
		return apierror.ErrorEmailIsRequired
	}
	if req.Password == "" {
		return apierror.ErrorPasswordIsRequired
	}
	if req.RepeatedPassword == "" {
		return apierror.ErrorRepeatedPasswordIsRequired
	}
	usernameLength := utf8.RuneCountInString(req.Username)
	if usernameLength < 2 || usernameLength > 32 {
		return apierror.ErrorInvalidUsernameLength
	}
	emailLength := utf8.RuneCountInString(req.Email)
	if emailLength < 2 || emailLength > 84 {
		return apierror.ErrorInvalidEmailLength
	}
	passwordLength := utf8.RuneCountInString(req.Password)
	if passwordLength < 8 || passwordLength > 64 {
		return apierror.ErrorInvalidPasswordLength
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, char := range req.Password {
		switch {
		case char > unicode.MaxASCII:
			return apierror.ErrorInvalidPasswordCharacter
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return apierror.ErrorMissingUpperInPassword
	}
	if !hasLower {
		return apierror.ErrorMissingLowerInPassword
	}
	if !hasDigit {
		return apierror.ErrorMissingDigitInPassword
	}

	if req.Password != req.RepeatedPassword {
		return apierror.ErrorPasswordsNotEqual
	}

	return nil
}
