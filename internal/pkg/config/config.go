package config

// Config is the object that contains the programm configuration.
type Config struct {
	Version int `json:"-"`

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
}
