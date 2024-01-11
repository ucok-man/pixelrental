package api

import (
	"sync"

	"github.com/ucok-man/pixelrental/internal/config"
	"github.com/ucok-man/pixelrental/internal/httpclient"
	"github.com/ucok-man/pixelrental/internal/logging"
	"github.com/ucok-man/pixelrental/internal/mailer"
	"github.com/ucok-man/pixelrental/internal/repo"
)

const version = "1.0.0"

type Application struct {
	config     *config.Config
	logger     *logging.Logger
	repo       *repo.Services
	mailer     *mailer.Mailer
	httpclient *httpclient.HTTPClient
	wg         sync.WaitGroup
	ctxkey     struct {
		user string
	}
}

func New() *Application {
	logger := logging.New()
	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err, "failed config initialization", nil)
	}

	dbconn, err := config.OpenDB(cfg)
	if err != nil {
		logger.Fatal(err, "failed open db connection", nil)
	}
	logger.Info("database connection pool established", nil)

	app := &Application{
		logger: logger,
		config: cfg,
		repo:   repo.New(dbconn),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
		ctxkey: struct {
			user string
		}{
			user: "user",
		},
		httpclient: httpclient.New(cfg),
	}

	return app
}
