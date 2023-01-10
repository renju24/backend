package apiserver

import (
	"net/http"

	"github.com/armantarkhanian/jwt"
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/config"
	"github.com/rs/zerolog"
)

// APIServer is the main object of programm.
type APIServer struct {
	runMode string
	addr    string

	router         *gin.Engine
	logger         *zerolog.Logger
	config         *config.Config
	jwt            *jwt.EncodeDecoder
	centrifugeNode *centrifuge.Node

	// Dependecies.
	db Database
	ConfigReader
}

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

// Run runs the HTTP server.
func (a *APIServer) Run() error {
	if a.runMode == "prod" {
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
				a.logger.Fatal().Err(err).Send()
			}
		}()
		return a.router.RunTLS(a.addr, "./cert.crt", "./private.key")
	} else {
		return a.router.Run(a.addr)
	}
}

// APIServer should be a singleton, so make it global.
var singleton *APIServer

// NewAPIServer creates a singleton APIServer object.
func NewAPIServer(runMode string, db Database, router *gin.Engine, logger *zerolog.Logger, configReader ConfigReader) *APIServer {
	if singleton == nil {
		singleton = initApi(runMode, db, router, logger, configReader)
	}
	return singleton
}

func initApi(runMode string, db Database, router *gin.Engine, logger *zerolog.Logger, configReader ConfigReader) *APIServer {
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
		runMode:      runMode,
		addr:         ":8008",
		router:       router,
		logger:       logger,
		config:       config,
		jwt:          jwtEncodeDecoder,
		db:           db,
		ConfigReader: configReader,
	}

	if a.runMode == "prod" {
		a.addr = ":443"
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowWebSockets = true

	// Static assets directory.
	// a.router.Use(static.Serve("/", static.LocalFile("./internal/apiserver/front", true)))

	a.router.Use(
		sameSiteMiddleware(a),
		cors.New(corsConfig),
		loggerMiddleware(a),
	)

	a.router.NoRoute(func(c *gin.Context) {
		// c.AbortWithStatus(http.StatusNotFound)
		c.File("./internal/apiserver/front/index.html")
		// return
	})

	// Static assets directory.
	a.router.Static("/assets", "./assets")

	a.router.StaticFile("/asset-manifest.json", "./internal/apiserver/front/asset-manifest.json")
	a.router.StaticFile("/favicon.ico", "./internal/apiserver/front/favicon.ico")
	a.router.StaticFile("/index.html", "./internal/apiserver/front/index.html")
	a.router.StaticFile("/logo192.png", "./internal/apiserver/front/logo192.png")
	a.router.StaticFile("/logo512.png", "./internal/apiserver/front/logo512.png")
	a.router.StaticFile("/logo.ico", "./internal/apiserver/front/logo.ico")
	a.router.StaticFile("/logo.png", "./internal/apiserver/front/logo.png")
	a.router.StaticFile("/manifest.json", "./internal/apiserver/front/manifest.json")
	a.router.StaticFile("/robots.txt", "./internal/apiserver/front/robots.txt")
	a.router.Static("/static", "./internal/apiserver/front/static")

	a.router.GET("/logout", func(c *gin.Context) {
		c.SetCookie(
			a.config.Server.Token.Cookie.Name,
			"",
			-1,
			a.config.Server.Token.Cookie.Path,
			a.config.Server.Token.Cookie.Domain,
			false,
			false,
		)
		c.Redirect(http.StatusFound, "/")
	})

	// POST /api/v1/*
	apiRoutes := a.router.Group("/api/v1")
	{
		apiRoutes.POST("/sign_up", signUp(a))
		apiRoutes.POST("/sign_in", signIn(a))
		apiRoutes.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "PONG") })
		apiRoutes.GET("/oauth2/:platform/services", oauth2Services(a))
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
		WebsocketConfig: centrifuge.WebsocketConfig{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			UseWriteBufferPool: true,
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
