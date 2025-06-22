package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"zendo/data_fetcher/routes"
	"zendo/data_fetcher/services"
	"zendo/data_fetcher/utils"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() {
	var encoding string
	var encoderCfg zapcore.EncoderConfig
	// TODO: flip to json logging in prod
	encoderCfg = zap.NewDevelopmentEncoderConfig()
	encoding = "console"

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// TODO: Make this configurable via env var
	zapLevel := zap.InfoLevel

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       os.Getenv("ZENDO_ENV") == "dev",
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          encoding,
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]any{
			"pid": os.Getpid(),
		},
	}
	logger := zap.Must(config.Build())

	zap.ReplaceGlobals(logger)
}

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env file")
	}

	setupLogger()

	// setup router
	mux := http.NewServeMux()

	// setup any auth / cors / logging middleware

	// setup dependencies

	httpClient := utils.HttpClient{}

	electricSerice := services.ElectricitymapService{
		Http: &httpClient,
	}
	weatherService := services.OpenMeteoWeatherService{
		Http: &httpClient,
	}

	dataService := services.CouchDBDataService{
		Http: &httpClient,
	}

	// setup routes and inject dependencies
	dataRoutes := routes.DataRoutes{
		ElectricService: &electricSerice,
		WeatherService:  &weatherService,
		DataService:     &dataService,
	}

	// register routes
	mux.HandleFunc("/update", dataRoutes.GetLatest)
	mux.HandleFunc("/seed", dataRoutes.Seed24Hrs)

	// configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	zap.L().Info("Server starting", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
