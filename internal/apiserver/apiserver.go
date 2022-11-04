package apiserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/armantarkhanian/jwt"
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/config"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
)

// APIServer is the main object of programm.
type APIServer struct {
	router         *gin.Engine
	logger         *zerolog.Logger
	config         *config.Config
	jwt            *jwt.EncodeDecoder
	centrifugeNode *centrifuge.Node

	// Dependecies.
	db Database
	ConfigReader
}

// Run runs the HTTP server.
func (a *APIServer) Run(port, runMode string) error {
	if runMode == "prod" {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		return autotls.RunWithContext(ctx, a.router, a.config.Server.Token.Cookie.Domain)
	} else {
		return a.router.Run(port)
	}
}

// APIServer should be a singleton, so make it global.
var singleton *APIServer

// NewAPIServer creates a singleton APIServer object.
func NewAPIServer(db Database, router *gin.Engine, logger *zerolog.Logger, configReader ConfigReader) *APIServer {
	if singleton == nil {
		singleton = initApi(db, router, logger, configReader)
	}
	return singleton
}

func initApi(db Database, router *gin.Engine, logger *zerolog.Logger, configReader ConfigReader) *APIServer {
	// Read config from database.
	logger.Info().Msg("Reading config from configReader.")
	config, err := configReader.ReadConfig()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	logger.Info().Interface("config", config).Send()

	jwtEncodeDecoder, err := jwt.New(config.Server.Token.SigningKey)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	a := &APIServer{
		router:       router,
		logger:       logger,
		config:       config,
		jwt:          jwtEncodeDecoder,
		db:           db,
		ConfigReader: configReader,
	}

	a.router.Use(
		corsMiddleware(a),
		loggerMiddleware(a),
	)

	a.router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	})

	// POST /api/v1/*
	apiRoutes := a.router.Group("/api/v1")
	{
		apiRoutes.POST("/sign_up", signUp(a))
		apiRoutes.POST("/sign_in", signIn(a))
		apiRoutes.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "PONG") })
	}

	// POST /api/v1/oauth2/login/:platform
	oauthLoginRoutes := apiRoutes.Group("/oauth2/login/:platform")
	{
		oauthLoginRoutes.GET("/google", func(c *gin.Context) {
			oauthConfig, err := googleOauthConfig(a, c.Param("platform"))
			if err != nil {
				log.Fatalln(err)
			}
			authPage := oauthConfig.AuthCodeURL("state")
			c.Redirect(http.StatusMovedPermanently, authPage)
		})
		oauthLoginRoutes.GET("/yandex", func(c *gin.Context) {
			oauthConfig, err := yandexOauthConfig(a, c.Param("platform"))
			if err != nil {
				log.Fatalln(err)
			}
			authPage := oauthConfig.AuthCodeURL("state")
			c.Redirect(http.StatusMovedPermanently, authPage)
		})
	}

	oauthCallbackRoutes := apiRoutes.Group("/oauth2/callback/:platform")
	{
		oauthCallbackRoutes.GET("/google", googleOauth(a))
		oauthCallbackRoutes.GET("/yandex", yandexOauth(a))
	}

	// Initialize WebSocket server.
	node, handler, err := websocket.New(websocket.Config{
		Engine:        &websocket.MemoryEngine{},
		ClientHandler: a,
		NodeHandler:   a,
		TokenLookup: websocket.TokenLookup{
			Header:       config.Server.Token.Header.Name,
			Cookie:       config.Server.Token.Cookie.Name,
			HeaderPrefix: "Bearer",
		},
	})
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	a.centrifugeNode = node

	// GET /connection/websocket
	a.router.GET("/connection/websocket", gin.WrapH(handler))

	return a
}

func googleOauthConfig(a *APIServer, platform string) (*oauth2.Config, error) {
	cfg := &oauth2.Config{
		ClientID:     a.config.Oauth2.Google.ClientID,
		ClientSecret: a.config.Oauth2.Google.ClientSecret,
		Scopes:       a.config.Oauth2.Google.Scopes,
		Endpoint:     google.Endpoint,
	}
	switch platform {
	case "web":
		cfg.RedirectURL = a.config.Oauth2.Google.Callbacks.Web
	case "android":
		cfg.RedirectURL = a.config.Oauth2.Google.Callbacks.Android
	default:
		return nil, errors.New("invalid platform")
	}
	return cfg, nil
}

func yandexOauthConfig(a *APIServer, platform string) (*oauth2.Config, error) {
	cfg := &oauth2.Config{
		ClientID:     a.config.Oauth2.Yandex.ClientID,
		ClientSecret: a.config.Oauth2.Yandex.ClientSecret,
		Scopes:       a.config.Oauth2.Yandex.Scopes,
		Endpoint:     yandex.Endpoint,
	}
	switch platform {
	case "web":
		cfg.RedirectURL = a.config.Oauth2.Yandex.Callbacks.Web
	case "android":
		cfg.RedirectURL = a.config.Oauth2.Yandex.Callbacks.Android
	default:
		return nil, errors.New("invalid platform")
	}
	return cfg, nil
}
