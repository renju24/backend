package apiserver

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/armantarkhanian/jwt"
	"github.com/armantarkhanian/websocket"
	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
	"github.com/renju24/backend/internal/pkg/apierror"
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

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

// Run runs the HTTP server.
func (a *APIServer) Run(port, runMode string) error {
	if runMode == "prod" {
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
				a.logger.Fatal().Err(err).Send()
			}
		}()

		return a.router.RunTLS(":443", "./cert.crt", "./private.key")
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

	// TODO: just for test. Remove it in production.
	a.router.SetHTMLTemplate(template.Must(template.New("").Parse(`
	{{ define "unauthorized.html" }}
		<b>Вы не авторизованы!</b><br><br>

		<a href="/api/v1/oauth2/web/google" style="outline: none;text-decoration: none;">
			<img src="https://cdn.iconscout.com/icon/free/png-256/google-1772223-1507807.png" width="30px" style="margin:10px;">
		</a>
		<a href="/api/v1/oauth2/web/yandex" style="outline: none;text-decoration: none;">
			<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/5/58/Yandex_icon.svg/2048px-Yandex_icon.svg.png" width="30px" style="margin:10px;">
		</a>
		<a href="/api/v1/oauth2/web/vk" style="outline: none;text-decoration: none;">
			<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/f3/VK_Compact_Logo_%282021-present%29.svg/2048px-VK_Compact_Logo_%282021-present%29.svg.png" width="30px" style="margin:10px;">
		</a>
	{{ end }}

	{{ define "authorized.html" }}
		<b>Вы авторизованы!</b><br><br>

		<b>ID</b>: {{ .ID }}<br>
		<b>Username</b>: {{ .Username }}<br>
		<b>Email</b>: {{ .Email }}<br><br>

		<a href="/logout">Логаут</a>
	{{ end }}
	`)))

	// Main page.
	a.router.GET("/", func(c *gin.Context) {
		token, _ := c.Cookie(a.config.Server.Token.Cookie.Name)
		var payload jwt.Payload
		if err := a.jwt.Decode(token, &payload); err != nil { // If user is not authorized.
			c.HTML(http.StatusOK, "unauthorized.html", nil)
			return
		}
		userID, err := strconv.ParseInt(payload.Subject, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, &apierror.Error{
				Error: apierror.ErrorInternal,
			})
			return
		}
		user, err := a.db.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, &apierror.Error{
				Error: apierror.ErrorInternal,
			})
			return
		}
		c.HTML(http.StatusOK, "authorized.html", map[string]interface{}{
			"ID":       user.ID,
			"Username": user.Username,
			"Email":    user.Email,
		})
	})

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
		apiRoutes.GET("/oauth2/services", oauth2Services(a))
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
