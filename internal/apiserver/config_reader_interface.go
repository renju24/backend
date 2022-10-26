package apiserver

import "github.com/renju24/backend/internal/pkg/config"

// ConfigReader is used to read config from multiple sources.
type ConfigReader interface {
	ReadConfig() (*config.Config, error)
}
