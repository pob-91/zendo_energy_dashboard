package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"zendo/api/routes"
	"zendo/lib_zendo/services"
	"zendo/lib_zendo/utils"

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
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// setup dependencies

	httpClient := utils.HttpClient{}

	dataService := services.CouchDBDataService{
		Http: &httpClient,
	}

	// setup routes and inject dependencies
	dataRoutes := routes.DataRoutes{
		DataService: &dataService,
	}

	// register routes
	mux.HandleFunc("/energy-summary", dataRoutes.GetLatestMetric)
	mux.HandleFunc("/historical-data", dataRoutes.GetTimeSeriesMetrics)

	// configure server
	server := &http.Server{
		Addr:         ":8081",
		Handler:      corsHandler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	zap.L().Info("Server starting", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
