package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/xlab/closer"
	"go.uber.org/zap"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

var sugar zap.SugaredLogger

func main() {

	parseFlags()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	mem := s.NewStorage()
	h.NewMemRepo(&mem)

	if FileStoragePath != "" {
		if Restore {
			err = h.ReadFileStorage(FileStoragePath)
			if err != nil {
				sugar.Infof("file loading error(%v)", err)
			}
		}
		if StoreInterval != 0 {
			go fileStorageRoutine()
		}

		closer.Bind(fileStorageExit)
		defer closer.Close()
	}

	r := chi.NewRouter()
	r.Post("/update/{memtype}/{name}/{value}", WithLogging(h.UpdateMem()))
	r.Get("/value/{memtype}/{name}", WithLogging(h.GetMem()))
	r.Post("/update/", WithLogging(gzipMiddleware(h.UpdateMemJSON())))
	r.Post("/value/", WithLogging(gzipMiddleware(h.GetMemJSON())))
	r.Get("/", WithLogging(gzipMiddleware(h.GetAllMems())))

	sugar.Infow("Starting server", "addr", flagAddr)

	err = http.ListenAndServe(flagAddr, r)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		panic(err)
	}
}
