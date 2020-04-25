package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/database"
	"github.com/pyr-sh/keybase-notarybot/bot/storage"
)

type Config struct {
	Addr     string
	Debug    bool
	HMACKey  []byte
	Username string

	Log      *zap.Logger
	Storage  *storage.Storage
	Database *database.Database
	Alice    *alice.Client
}

type API struct {
	Config
	engine *gin.Engine
	server *http.Server
}

func New(cfg Config) (*API, error) {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	api := &API{
		Config: cfg,
		engine: engine,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: cors.Default().Handler(engine),
		},
	}
	if err := api.Routes(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *API) Routes() error {
	a.engine.Use(a.errorMiddleware)

	a.engine.POST("/signatures", a.signaturesCreate)
	return nil
}

func (a *API) Start(ctx context.Context) error {
	a.Log.With(zap.String("addr", a.Addr)).Info("Starting the API server")
	if err := a.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
