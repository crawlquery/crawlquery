package testutil

import "go.uber.org/zap"

func NewTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
