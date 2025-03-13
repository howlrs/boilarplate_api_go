package main

import (
	"backend/routes"
	"fmt"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

var (
	ADDRESS string
	isTest  bool
)

func init() {
	// 環境別の処理 (Environment-specific processing)
	isProduction := runtime.GOOS == "linux"

	// Load appropriate .env file
	if isProduction {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load(".env.local")
	}

	// Configure test mode
	isTest = !isProduction || os.Getenv("ISTEST") == "true"

	// Configure logging level
	configureLogging(isProduction, isTest)

	// Configure PORT
	port := os.Getenv("PORT")
	if isProduction {
		ADDRESS = fmt.Sprintf("0.0.0.0:%s", port)
		log.Debug().Msgf("Linuxでの処理, ADDRESS: %s", ADDRESS)
	} else {
		ADDRESS = fmt.Sprintf("localhost:%s", port)
		log.Debug().Msgf("その他のOSでの処理, ADDRESS: %s", ADDRESS)
	}
}

func main() {
	e := echo.New()
	routes.Endpoint(e, isTest)

	log.Fatal().Err(e.Start(ADDRESS))
}

func configureLogging(isProduction bool, isTest bool) {
	if isProduction && !isTest {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if isTest {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("テストモードで起動します")
	} else {
		// Development environment with custom log level
		logLevel := zerolog.ErrorLevel
		switch os.Getenv("LOG_LEVEL") {
		case "debug":
			logLevel = zerolog.DebugLevel
		case "info":
			logLevel = zerolog.InfoLevel
		case "warn":
			logLevel = zerolog.WarnLevel
		case "error":
			logLevel = zerolog.ErrorLevel
		}
		zerolog.SetGlobalLevel(logLevel)
	}
}
