package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nstapelbroek/gatekeeper/app/adapters"
	"github.com/nstapelbroek/gatekeeper/app/handlers"
	"github.com/nstapelbroek/gatekeeper/app/middlewares"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// App is a struct holding references to every gatekeeper component
type App struct {
	router            *gin.Engine
	config            *viper.Viper
	adapterFactory    *adapters.AdapterFactory
	adapterDispatcher *adapters.AdapterDispatcher
	logger            *zap.Logger
	//register       *domain.Register
}

// NewApp is a constructor method for App
func NewApp(c *viper.Viper) *App {
	a := App{
		config: c,
	}

	bootRouter(&a)
	bootLogging(&a)
	bootServices(&a)
	bootMiddleware(&a)
	bootRoutes(&a)

	return &a
}

func bootLogging(a *App) {
	logLevelConfig := zap.NewAtomicLevel()
	if gin.Mode() == gin.DebugMode {
		logLevelConfig.SetLevel(zapcore.DebugLevel)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		logLevelConfig,
	))
	a.logger = logger
}

func bootMiddleware(a *App) {
	middlewares.RegisterAccessLogMiddleware(a.router, a.logger)
	middlewares.RegisterBasicAuthentication(a.router, a.config)
}

func bootServices(a *App) {
	adapterFactory, err := adapters.NewAdapterFactory(a.config)
	a.adapterFactory = adapterFactory
	if err != nil {
		log.Fatalln(err)
	}

	dispatcher, err := adapters.NewAdapterDispatcher(a.adapterFactory.GetAdapters(), a.logger)
	a.adapterDispatcher = dispatcher
	if err != nil {
		log.Fatalln(err)
	}
}

func bootRouter(a *App) {
	gin.SetMode(a.config.GetString("app_env"))
	a.router = gin.New()
	a.router.HandleMethodNotAllowed = true
}

func bootRoutes(a *App) {
	gateHandler, err := handlers.NewGateHandler(
		a.config.GetInt64("rule_close_timeout"),
		a.config.GetString("rule_ports"),
		a.adapterDispatcher,
		a.logger,
	)

	if err != nil {
		panic(err)
	}

	a.router.POST("/", gateHandler.PostOpen)
	a.router.NoRoute(handlers.NotFound)
	a.router.NoMethod(handlers.MethodNotAllowed)
}

// Run will start the HTTP server of the app instance
func (a App) Run() (err error) {
	port := a.config.GetInt("http_port")
	a.logger.Info("Starting gatekeeper", zap.Int("port", port))
	err = a.router.Run(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Println(err.Error())
	}

	_ = a.logger.Sync()

	return
}
