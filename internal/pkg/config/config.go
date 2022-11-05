package config

// Config is the object that contains the programm configuration.
type Config struct {
	Version int `json:"version"`

	Server struct {
		CORS struct {
			SameSite                      string `json:"same_site"`
			AccessControlAllowOrigin      string `json:"access_control_allow_origin"`
			AccessControlAllowCredentials bool   `json:"access_control_allow_credentials"`
		} `json:"cors"`

		Token struct {
			SigningKey string `json:"signing_key"`
			Cookie     struct {
				Name     string `json:"name"`
				MaxAge   int    `json:"max_age"`
				Path     string `json:"path"`
				Domain   string `json:"domain"`
				Secure   bool   `json:"secure"`
				HttpOnly bool   `json:"http_only"`
			} `json:"cookie"`
			Header struct {
				Name string `json:"name"`
			} `json:"header"`
		} `json:"token"`
	} `json:"server"`

	Oauth2 OauthConfig `json:"oauth2"`
}

type OauthConfig struct {
	DeepLinks OauthRedirects      `json:"deep_links"`
	Google    OauthProviderConfig `json:"google"`
	Yandex    OauthProviderConfig `json:"yandex"`
	Github    OauthProviderConfig `json:"github"`
	VK        OauthProviderConfig `json:"vk"`
}

type OauthProviderConfig struct {
	ClientID     string         `json:"client_id"`
	ClientSecret string         `json:"client_secret"`
	Scopes       []string       `json:"scopes"`
	Callbacks    OauthRedirects `json:"callbacks"`
	API          string         `json:"api_url"`
}

type OauthRedirects struct {
	Web     string `json:"web"`
	Android string `json:"android"`
	IOS     string `json:"ios"`
}
