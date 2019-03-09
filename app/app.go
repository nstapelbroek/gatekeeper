package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
	"github.com/nstapelbroek/gatekeeper/app/handlers"
	"github.com/nstapelbroek/gatekeeper/app/middlewares"
	"github.com/spf13/viper"
)

type App struct {
	router         *gin.Engine
	config         *viper.Viper
	adapterFactory *adapters.AdapterFactory
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
	a.adapterFactory = adapters.NewAdapterFactory(a.config)
}

func (a App) Run() error {
	return a.router.Run(":" + a.config.GetString("http_port"))
}

func bootRouter(a *App) {
	gin.SetMode(a.config.GetString("app_env"))
	a.router = gin.Default()
}

func bootRoutes(a *App) {
	var adapterSlice []adapters.Adapter
	adapterSlice = append(adapterSlice, a.adapterFactory.GetAdapter())
	gateHandler, err := handlers.NewGateHandler(
		a.config.GetInt64("RULE_CLOSE_TIMEOUT"),
		a.config.GetString("RULE_PORTS"),
		adapterSlice,
	)

	if err != nil {
		panic(err)
	}

	a.router.POST("/", gateHandler.PostOpen)
	a.router.Handle("GET", "/", handlers.MethodNotAllowed)
	a.router.Handle("PATCH", "/", handlers.MethodNotAllowed)
	a.router.Handle("PUT", "/", handlers.MethodNotAllowed)
	a.router.Handle("HEAD", "/", handlers.MethodNotAllowed)
	a.router.Handle("OPTIONS", "/", handlers.MethodNotAllowed)
	a.router.NoRoute(handlers.NotFound)
}
