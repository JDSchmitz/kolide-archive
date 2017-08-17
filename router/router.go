package router

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/JDSchmitz/kolide-archive/config"
	"github.com/JDSchmitz/kolide-archive/controller"
	"github.com/JDSchmitz/kolide-archive/router/middleware/gzip"
	"github.com/JDSchmitz/kolide-archive/router/middleware/header"
	"github.com/JDSchmitz/kolide-archive/router/middleware/location"
	"github.com/JDSchmitz/kolide-archive/router/middleware/requestlogger"
	"github.com/JDSchmitz/kolide-archive/router/middleware/session"
	"github.com/JDSchmitz/kolide-archive/router/v1"
	"github.com/JDSchmitz/kolide-archive/shared/hub"
	"github.com/JDSchmitz/kolide-archive/static"
	"github.com/JDSchmitz/kolide-archive/version"

	"github.com/gin-gonic/contrib/expvar"
	"github.com/gin-gonic/contrib/sessions"
	// "github.com/mephux/contrib/expvar"
	// "github.com/mephux/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Load will setup and configure the gin router
func Load(configuration *config.Config) http.Handler {
	e := gin.New()

	e.SetHTMLTemplate(static.Load(configuration))

	// e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(gzip.Gzip(gzip.DefaultCompression))
	e.Use(session.SetUser())

	e.Use(location.Resolve)
	e.Use(header.NoCache)
	e.Use(header.Options)

	e.Use(header.Secure)

	e.Use(requestlogger.New(log.StandardLogger(), time.RFC3339, false))

	// if configuration.Session.Type == "cookie" {
	// store := sessions.NewCookieStore([]byte(configuration.Session.Key))
	// e.Use(sessions.Sessions(configuration.Session.Name, store))

	// } else if configuration.Session.Type == "redis" {
	store, err := sessions.NewRedisStore(configuration.Session.Size,
		configuration.Session.Network, configuration.Session.Address,
		configuration.Session.Password, []byte(configuration.Session.Key))

	if err != nil {
		log.Fatal(err)
	}

	e.Use(sessions.Sessions(version.Name, store))
	// } else {
	// log.Fatalf("unknown session type: %s", configuration.Session.Type)
	// }

	if configuration.Server.Production {
		log.Info("Loading Assets: MEMORY")
		e.StaticFS("/assets/", static.FileSystem())
	} else {
		log.Info("Loading Assets: DISK")
		e.StaticFS("/assets", http.Dir("./static/ui/"))

		var debugRoute = "/debug"
		log.Infof("Setting debug route: %s\n", debugRoute)
		e.GET(debugRoute, expvar.Handler())
	}

	// 404
	e.NoRoute(controller.Error)
	e.GET("/", controller.Index)
	e.GET("/login", controller.Login)

	// register v1 routes
	v1.Register(e, configuration)

	// websockets
	go hub.Websocket.Run()

	return e
}
