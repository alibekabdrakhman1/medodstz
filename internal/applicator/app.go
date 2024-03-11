package applicator

import (
	"context"
	"github.com/alibekabdrakhman1/medodsTZ/internal/config"
	"github.com/alibekabdrakhman1/medodsTZ/internal/controller"
	"github.com/alibekabdrakhman1/medodsTZ/internal/controller/http"
	"github.com/alibekabdrakhman1/medodsTZ/internal/service"
	"github.com/alibekabdrakhman1/medodsTZ/internal/storage"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

type App struct {
	logger *zap.SugaredLogger
	config *config.Config
}

func New(logger *zap.SugaredLogger, config *config.Config) *App {
	return &App{
		logger: logger,
		config: config,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	gracefullyShutdown(cancel)
	repo, err := storage.NewRepository(ctx, a.config, a.logger)
	if err != nil {
		log.Fatalf("cannot сonnect to db '%s:%s': %v", a.config.Database.Host, a.config.Database.Port, err)
	}
	authService := service.NewManager(repo, a.config.Auth.JwtSecretKey, a.logger)

	endPointHandler := http.NewManager(authService, a.logger)
	HTTPServer := controller.NewServer(a.config, endPointHandler)

	return HTTPServer.StartHTTPServer(ctx)
}

func gracefullyShutdown(c context.CancelFunc) {
	osC := make(chan os.Signal, 1)
	signal.Notify(osC, os.Interrupt)
	go func() {
		log.Print(<-osC)
		c()
	}()
}
