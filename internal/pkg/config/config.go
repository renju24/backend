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

	Oauth2 struct {
		DeepLinks OauthRedirects `json:"deep_links"`
		Google    OauthConfig    `json:"google"`
		Yandex    OauthConfig    `json:"yandex"`
		Github    OauthConfig    `json:"github"`
	} `json:"oauth2"`
}

type Platform string

const (
	Web     Platform = "web"
	Android Platform = "android"
	IOS     Platform = "ios"
)

type OauthService string

const (
	Google OauthService = "google"
	Yandex OauthService = "yandex"
	Github OauthService = "github"
)

type OauthConfig struct {
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
