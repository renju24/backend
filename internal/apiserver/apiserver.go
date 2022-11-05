package apiserver

import (
	"context"
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
		apiRoutes.GET("/oauth2/:platform/:service", oauth2Login(a))
		apiRoutes.GET("/oauth2/:platform/:service/callback", oauth2Callback(a))
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
