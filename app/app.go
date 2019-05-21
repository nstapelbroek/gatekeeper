package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
	"github.com/nstapelbroek/gatekeeper/app/handlers"
	"github.com/nstapelbroek/gatekeeper/app/middlewares"
	"github.com/spf13/viper"
	"log"
)

type App struct {
	router            *gin.Engine
	config            *viper.Viper
	adapterFactory    *adapters.AdapterFactory
	adapterDispatcher *adapters.AdapterDispatcher
	//register       *domain.Register
}

func NewApp(c *viper.Viper) *App {
	a := App{
		config: c,
	}

	bootServices(&a)
	bootRouter(&a)
	bootMiddleware(&a)
	bootRoutes(&a)

	return &a
}

func bootMiddleware(a *App) {
	middlewares.RegisterResolverMiddleware(a.router, a.config)
	middlewares.RegisterBasicAuthentication(a.router, a.config)
}

func bootServices(a *App) {
	adapterFactory, err := adapters.NewAdapterFactory(a.config)
	a.adapterFactory = adapterFactory
	if err != nil {
		log.Fatalln(err)
	}

	dispatcher, err := adapters.NewAdapterDispatcher(a.adapterFactory.GetAdapters())
	a.adapterDispatcher = dispatcher
	if err != nil {
		log.Fatalln(err)
	}
}

func bootRouter(a *App) {
	gin.SetMode(a.config.GetString("app_env"))
	a.router = gin.Default()
	a.router.HandleMethodNotAllowed = true
}

func bootRoutes(a *App) {
	gateHandler, err := handlers.NewGateHandler(
		a.config.GetInt64("rule_close_timeout"),
		a.config.GetString("rule_ports"),
		a.adapterDispatcher,
	)

	if err != nil {
		panic(err)
	}

	a.router.POST("/", gateHandler.PostOpen)
	a.router.NoRoute(handlers.NotFound)
	a.router.NoMethod(handlers.MethodNotAllowed)
}

func (a App) Run() (err error) {
	err = a.router.Run(":" + a.config.GetString("http_port"))
	if err != nil {
		log.Println(err.Error())
	}

	return
}
