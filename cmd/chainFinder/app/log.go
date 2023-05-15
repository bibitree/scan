package app

import (
	"ethgo/chainFinder"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func SetLogger(logger *zap.SugaredLogger) {
	log = logger
	chainFinder.SetLogger(logger)
}
