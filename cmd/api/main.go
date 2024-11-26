package main

import (
	"log"
	"log/slog"
	"os"

	"cyansnbrst.com/order-info/config"
	"cyansnbrst.com/order-info/internal/server"
	"cyansnbrst.com/order-info/pkg/db/postgres"
)

func main() {
	log.Println("starting api server")

	cfgFile, err := config.LoadConfig("config/config-local")
	if err != nil {
		log.Fatalf("loadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("parseConfig: %v", err)
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	psqlDB, err := postgres.OpenDB(cfg)
	if err != nil {
		logger.Error("failed to init storage:", "error", err)
		os.Exit(1)
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, logger, psqlDB)
	if err = s.Run(); err != nil {
		logger.Error("an error occured")
		os.Exit(1)
	}
}
