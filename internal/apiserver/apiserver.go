package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/apierror"
	"github.com/renju24/backend/internal/pkg/config"
	"github.com/rs/zerolog"
)

// APIError is the JSON-object that server will return when an error occurs.
type APIError struct {
	Status int             `json:"status"`
	Error  *apierror.Error `json:"error"`
}

// APIServer is the main object of programm.
type APIServer struct {
	router *gin.Engine
	logger *zerolog.Logger
	config *config.Config

	// Dependecies.
	db Database
}

// Run runs the HTTP server.
func (a *APIServer) Run() error {
	return a.router.Run(a.config.Server.Addr)
}

// APIServer should be a singleton, so make it global.
var singleton *APIServer

// NewAPIServer creates a singleton APIServer object.
func NewAPIServer(db Database, router *gin.Engine, logger *zerolog.Logger, config *config.Config) *APIServer {
	if singleton == nil {
		singleton = initApi(db, router, logger, config)
	}
	return singleton
}

func initApi(db Database, router *gin.Engine, logger *zerolog.Logger, config *config.Config) *APIServer {
	a := &APIServer{
		router: router,
		logger: logger,
		config: config,
		db:     db,
	}

	a.router.Use(
		corsMiddleware(a),
		loggerMiddleware(a),
	)

	a.router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	})

	apiRoutes := a.router.Group("/api/v1")
	{
		apiRoutes.POST("/api/v1/sign_up", signUp(a))
	}

	return a
}
