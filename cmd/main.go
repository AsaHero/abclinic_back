package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AsaHero/abclinic/internal/app"
	"github.com/AsaHero/abclinic/internal/pkg/config"
	"go.uber.org/zap"
)

func main() {
	// load dot env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("cannot Load .env: %v", err)
	// }

	// config init
	cfg := config.NewConfig()

	// app init
	app := app.NewApp(cfg)

	// app runs
	go func() {
		app.Logger.Info("Listen:", zap.String("address", cfg.Server.Host+cfg.Server.Port))
		if err := app.Run(); err != nil {
			app.Logger.Error("error while running server: %v", zap.Error(err))
		}
	}()

	// wait for sigint
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	app.Logger.Info("abclinic sevrver stops")
	app.Stop()
}
