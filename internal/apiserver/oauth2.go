package apiserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
)

// TODO: hard code.
func oauth2Services(api *APIServer) gin.HandlerFunc {
	type imageAndURL struct {
		Image string `json:"image"`
		URL   string `json:"url"`
	}
	type service struct {
		Name    config.OauthService `json:"name"`
		Web     imageAndURL         `json:"web"`
		Android imageAndURL         `json:"android"`
	}
	type response struct {
		Services []service `json:"services"`
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &response{
			Services: []service{
				{
					Name: config.Google,
					Web: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Google.Callbacks.Web, "/callback"),
					},
					Android: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Google.Callbacks.Android, "/callback"),
					},
				},
				{
					Name: config.Yandex,
					Web: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Yandex.Callbacks.Web, "/callback"),
					},
					Android: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Yandex.Callbacks.Android, "/callback"),
					},
				},
				{
					Name: config.Github,
					Web: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Github.Callbacks.Web, "/callback"),
					},
					Android: imageAndURL{
						Image: "",
						URL:   strings.TrimSuffix(api.config.Oauth2.Github.Callbacks.Android, "/callback"),
					},
				},
			},
		})
	}
}

func oauth2Login(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := parseService(c.Param("service"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		platform, err := parsePlatform(c.Param("platform"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		var oauthCfg *oauth2.Config
		switch service {
		case config.Google:
			oauthCfg, err = oauthConfig(api, config.Google, platform)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		case config.Yandex:
			oauthCfg, err = oauthConfig(api, config.Yandex, platform)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		case config.Github:
			oauthCfg, err = oauthConfig(api, config.Github, platform)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		default:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		authPage := oauthCfg.AuthCodeURL("state")
		c.Redirect(http.StatusMovedPermanently, authPage)
	}
}

func oauth2Callback(api *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := parseService(c.Param("service"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		platform, err := parsePlatform(c.Param("platform"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		switch service {
		case config.Google:
			googleOauth(api, c, service, platform)
		case config.Yandex:
			yandexOauth(api, c, service, platform)
		case config.Github:
			githubOauth(api, c, service, platform)
		default:
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

var (
	ErrUnknownPlatform = errors.New("unknown platform")
	ErrUnknownService  = errors.New("unknown service")
)

func parsePlatform(s string) (config.Platform, error) {
	switch s {
	case "web":
		return config.Web, nil
	case "android":
		return config.Android, nil
	}
	return "", ErrUnknownPlatform
}

func parseService(s string) (config.OauthService, error) {
	switch s {
	case "google":
		return config.Google, nil
	case "yandex":
		return config.Yandex, nil
	case "github":
		return config.Github, nil
	}
	return "", ErrUnknownService
}

func oauthConfig(a *APIServer, service config.OauthService, platform config.Platform) (*oauth2.Config, error) {
	switch service {
	case config.Google:
		cfg := &oauth2.Config{
			ClientID:     a.config.Oauth2.Google.ClientID,
			ClientSecret: a.config.Oauth2.Google.ClientSecret,
			Scopes:       a.config.Oauth2.Google.Scopes,
			Endpoint:     google.Endpoint,
		}
		switch platform {
		case config.Web:
			cfg.RedirectURL = a.config.Oauth2.Google.Callbacks.Web
		case config.Android:
			cfg.RedirectURL = a.config.Oauth2.Google.Callbacks.Android
		}
		return cfg, nil
	case config.Yandex:
		cfg := &oauth2.Config{
			ClientID:     a.config.Oauth2.Yandex.ClientID,
			ClientSecret: a.config.Oauth2.Yandex.ClientSecret,
			Scopes:       a.config.Oauth2.Yandex.Scopes,
			Endpoint:     yandex.Endpoint,
		}
		switch platform {
		case config.Web:
			cfg.RedirectURL = a.config.Oauth2.Yandex.Callbacks.Web
		case config.Android:
			cfg.RedirectURL = a.config.Oauth2.Yandex.Callbacks.Android
		}
		return cfg, nil
	case config.Github:
		cfg := &oauth2.Config{
			ClientID:     a.config.Oauth2.Github.ClientID,
			ClientSecret: a.config.Oauth2.Github.ClientSecret,
			Scopes:       a.config.Oauth2.Github.Scopes,
			Endpoint:     github.Endpoint,
		}
		switch platform {
		case config.Web:
			cfg.RedirectURL = a.config.Oauth2.Github.Callbacks.Web
		case config.Android:
			cfg.RedirectURL = a.config.Oauth2.Github.Callbacks.Android
		}
		return cfg, nil
	}
	return nil, ErrUnknownService
}
