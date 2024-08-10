package internal

import (
	"go.uber.org/zap"
)

var Logger zap.SugaredLogger

func InitLogger() {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:      "json",
		EncoderConfig: zap.NewDevelopmentEncoderConfig(),
		OutputPaths:   []string{"debug.log"},
	}

	//logger, err := zap.NewDevelopment()
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
	}(logger)

	if err != nil {
		panic(err)
	}
	Logger = *logger.Sugar()
}
