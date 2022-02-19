package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Init() {
	var config zap.Config
	if os.Getenv("DEV") == "true" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()

		if os.Getenv("LOGS_PATH") != "" {
			config.OutputPaths = []string{os.Getenv("LOGS_PATH")}
		}
	}
	// config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)

	if godotenv.Load() != nil {
		zap.S().Fatal("Error loading .env file")
	}
}