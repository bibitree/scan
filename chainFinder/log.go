package chainFinder

import (
	"ethgo/util/logx"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger = logx.Default("[ChainFinder]")

func SetLogger(logger *zap.SugaredLogger) {
	log = logger
}
