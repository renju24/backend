package oauth2

import (
	"github.com/renju24/backend/internal/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/vk"
	"golang.org/x/oauth2/yandex"
)

type User struct {
	ID       string
	Username string
	Email    *string
}

func OauthConfig(providers config.OauthConfig, service Service, platform Platform) (*oauth2.Config, error) {
	switch service {
	case Google:
		cfg := &oauth2.Config{
			ClientID:     providers.Google.ClientID,
			ClientSecret: providers.Google.ClientSecret,
			Scopes:       providers.Google.Scopes,
			Endpoint:     google.Endpoint,
		}
		switch platform {
		case Web:
			cfg.RedirectURL = providers.Google.Callbacks.Web
		case Android:
			cfg.RedirectURL = providers.Google.Callbacks.Android
		}
		return cfg, nil
	case Yandex:
		cfg := &oauth2.Config{
			ClientID:     providers.Yandex.ClientID,
			ClientSecret: providers.Yandex.ClientSecret,
			Scopes:       providers.Yandex.Scopes,
			Endpoint:     yandex.Endpoint,
		}
		switch platform {
		case Web:
			cfg.RedirectURL = providers.Yandex.Callbacks.Web
		case Android:
			cfg.RedirectURL = providers.Yandex.Callbacks.Android
		}
		return cfg, nil
	case Github:
		cfg := &oauth2.Config{
			ClientID:     providers.Github.ClientID,
			ClientSecret: providers.Github.ClientSecret,
			Scopes:       providers.Github.Scopes,
			Endpoint:     github.Endpoint,
		}
		switch platform {
		case Web:
			cfg.RedirectURL = providers.Github.Callbacks.Web
		case Android:
			cfg.RedirectURL = providers.Github.Callbacks.Android
		}
		return cfg, nil
	case VK:
		cfg := &oauth2.Config{
			ClientID:     providers.VK.ClientID,
			ClientSecret: providers.VK.ClientSecret,
			Scopes:       providers.VK.Scopes,
			Endpoint:     vk.Endpoint,
		}
		switch platform {
		case Web:
			cfg.RedirectURL = providers.VK.Callbacks.Web
		case Android:
			cfg.RedirectURL = providers.VK.Callbacks.Android
		}
		return cfg, nil
	}
	return nil, ErrUnknownService
}
