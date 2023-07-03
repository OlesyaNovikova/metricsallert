package main

import (
	"time"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
)

func fileStorageRoutine() {

	ticker := time.NewTicker(StoreInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := h.WriteFileStorage(FileStoragePath)
		if err != nil {
			sugar.Errorf("file write error(%v)", err)
			return
		}
	}
}

func fileStorageExit() {
	err := h.WriteFileStorage(FileStoragePath)
	if err != nil {
		sugar.Errorf("file write error(%v)", err)
		return
	}
}
