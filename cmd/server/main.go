package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/xlab/closer"
	"go.uber.org/zap"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
	m "github.com/OlesyaNovikova/metricsallert.git/internal/middleware"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

func main() {

	parseFlags()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()

	mem, err := s.NewFileStorage(FileStoragePath, Restore, StoreInterval)
	if err != nil {
		sugar.Error(err.Error())
	} else {
		closer.Bind(mem.FileStorageExit)
		defer closer.Close()
	}
	h.NewMemRepo(mem)

	r := chi.NewRouter()
	r.Post("/update/{memtype}/{name}/{value}", m.WithLogging(sugar, h.UpdateMem()))
	r.Get("/value/{memtype}/{name}", m.WithLogging(sugar, h.GetMem()))
	r.Post("/update/", m.WithLogging(sugar, m.GzipMiddleware(h.UpdateMemJSON())))
	r.Post("/value/", m.WithLogging(sugar, m.GzipMiddleware(h.GetMemJSON())))
	r.Get("/", m.WithLogging(sugar, m.GzipMiddleware(h.GetAllMems())))

	sugar.Infow("Starting server", "addr", flagAddr)

	err = http.ListenAndServe(flagAddr, r)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		panic(err)
	}
}
