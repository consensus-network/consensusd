package main

import (
	"github.com/consensus-network/consensusd/infrastructure/logger"
	"github.com/consensus-network/consensusd/util/panics"
)

var (
	backendLog = logger.NewBackend()
	log        = backendLog.Logger("RPIC")
	spawn      = panics.GoroutineWrapperFunc(log)
)
