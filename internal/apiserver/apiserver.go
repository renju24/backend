package apiserver

import (
	"context"
	"encoding/json"
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

// APIError is the JSON-object that server will return when an error occurs.
type APIError struct {
	Error *centrifuge.Error `json:"error"`
}

type apiErrorJSON struct {
	Error errorJSON `json:"error"`
}

type errorJSON struct {
	Code      uint32 `json:"code"`
	Message   string `json:"message"`
	Temporary bool   `json:"temporary"`
}

// MarshalJSON ...
func (e *APIError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&apiErrorJSON{
		Error: errorJSON{
			Code:      e.Error.Code,
			Message:   e.Error.Message,
			Temporary: e.Error.Temporary,
		},
	})
}

// MarshalJSON ...
func (e *APIError) UnmarshalJSON(data []byte) error {
	var apiErrorJSON apiErrorJSON
	if err := json.Unmarshal(data, &apiErrorJSON); err != nil {
		return err
	}
	e.Error = &centrifuge.Error{
		Code:      apiErrorJSON.Error.Code,
		Message:   apiErrorJSON.Error.Message,
		Temporary: apiErrorJSON.Error.Temporary,
	}
	return nil
}

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
