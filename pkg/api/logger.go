package api

import (
	"go.uber.org/zap"
	"sync"
)

var lock = &sync.Mutex{}

var singleInstance *zap.Logger

func GetLogger() *zap.Logger {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance, _ = zap.NewProduction()
		}
	}

	return singleInstance
}
