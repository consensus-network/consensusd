package rpcclient

import (
	"github.com/consensus-network/consensusd/infrastructure/logger"
	"github.com/consensus-network/consensusd/util/panics"
)

var log = logger.RegisterSubSystem("RPCC")
var spawn = panics.GoroutineWrapperFunc(log)
