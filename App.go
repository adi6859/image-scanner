package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/devtron-labs/image-scanner/internal/middleware"
	"net/http"
	"os"
	"time"

	client "github.com/devtron-labs/common-lib/pubsub-lib"
	"github.com/devtron-labs/image-scanner/api"
	"github.com/devtron-labs/image-scanner/pubsub"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type App struct {
	Router           *api.Router
	Logger           *zap.SugaredLogger
	server           *http.Server
	db               *pg.DB
	natsSubscription *pubsub.NatSubscriptionImpl
	pubSubClient     *client.PubSubClientServiceImpl
}

func NewApp(Router *api.Router, Logger *zap.SugaredLogger,
	db *pg.DB, natsSubscription *pubsub.NatSubscriptionImpl,
	pubSubClient *client.PubSubClientServiceImpl) *App {
	return &App{
		Router:           Router,
		Logger:           Logger,
		db:               db,
		natsSubscription: natsSubscription,
		pubSubClient:     pubSubClient,
	}
}

type ServerConfig struct {
	SERVER_HTTP_PORT int `env:"SERVER_HTTP_PORT" envDefault:"8080"`
}

func (app *App) Start() {
	serverConfig := ServerConfig{}
	err := env.Parse(&serverConfig)
	if err != nil {
		app.Logger.Errorw("error in parsing server config from environment", "err", err)
		os.Exit(2)
	}
	httpPort := serverConfig.SERVER_HTTP_PORT
	app.Logger.Infow("starting server on ", "httpPort", httpPort)
	app.Router.Init()
	server := &http.Server{Addr: fmt.Sprintf(":%d", httpPort), Handler: app.Router.Router}
	app.Router.Router.Use(middleware.PrometheusMiddleware)
	app.server = server
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Errorw("error in startup", "err", err)
		os.Exit(2)
	}
}

func (app *App) Stop() {
	app.Logger.Infow("image scanner shutdown initiating")
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Logger.Infow("closing router")
	err := app.server.Shutdown(timeoutContext)
	if err != nil {
		app.Logger.Errorw("error in mux router shutdown", "err", err)
	}
	app.Logger.Infow("closing db connection")
	err = app.db.Close()
	if err != nil {
		app.Logger.Errorw("Error while closing DB", "error", err)
	}

	app.Logger.Infow("housekeeping done. exiting now")
}
