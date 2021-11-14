package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
	"webserver/config"
	"webserver/controller"
	"webserver/database"
	"webserver/logger"
	"webserver/server"
	"webserver/service"
	tracerCreator "webserver/tracing/creator"
	"webserver/transport"
	"webserver/utility"
)

var (
	configFile string
)

func parseParameters() {
	flag.StringVar(&configFile, "config", "", "the config file used to launch the service")
	flag.Parse()

	if len(configFile) == 0 {
		log.Fatalf("missing config file")
	}
}

func main() {
	parseParameters()

	configuration, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config failed with error: %v", err)
	}

	zapLogger, err := logger.SetupLogger(configuration.Logger)
	if err != nil {
		log.Fatalf("setup logger failed with error: %v", err)
	}

	tracer, err := tracerCreator.NewTracer("webserver", "", zapLogger)
	if err != nil {
		zapLogger.Fatal("create tracer failed",
			zap.Error(err))
	}
	defer tracer.Close()

	db, err := database.NewDatabase(zapLogger, configuration.DB)
	if err != nil {
		zapLogger.Fatal("create database failed",
			zap.Error(err))
	}
	defer db.Close()

	userService := service.NewUserService(zapLogger, db)
	userController := controller.NewUserController(zapLogger, userService)
	svr := server.NewServer(zapLogger)

	if err := svr.SetupServer(configuration.Server); err != nil {
		zapLogger.Info("setup server failed", zap.Error(err))
		os.Exit(1)
	}

	transport.SetupRouters(svr, userController)

	zapLogger.Info("start server")

	g := utility.MakeGroup()
	g.Add(svr.RunServer, svr.StopServer)

	if err := g.Run(); err != nil {
		zapLogger.Info("run failed", zap.Error(err))
		os.Exit(1)
	}
}
