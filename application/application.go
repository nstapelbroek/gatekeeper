package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/nstapelbroek/gatekeeper/adapters"
	"github.com/nstapelbroek/gatekeeper/middlewares"
	"github.com/nstapelbroek/gatekeeper/controllers"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	app := &Application{}
	app.config = config

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config *viper.Viper
}

// MiddlewareStruct is used for bootstrapping and loading the interpose middleware in mux
func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.MustAuthenticate(app.config))
	middle.Use(middlewares.ResolveOrigin(app.config))
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()
	handler := controllers.NewGateController(
		adapters.NewAdapterFactory(app.config),
		app.config.GetInt("closure_timeout"),
	)

	router.Handle("/", http.HandlerFunc(handler.PostOpen)).Methods("POST")

	// Due to the first-match approach of Gorilla mux, we serve the static files last
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
