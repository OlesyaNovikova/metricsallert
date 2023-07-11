package main

import (
	"context"
	"database/sql"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xlab/closer"
	"go.uber.org/zap"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
	m "github.com/OlesyaNovikova/metricsallert.git/internal/middleware"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

func main() {

	parseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()

	if DBAddr != "" {
		db, err := sql.Open("pgx", DBAddr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		ctx = context.WithValue(ctx, "db", db)
	}

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
	r.Get("/ping", m.WithLogging(sugar, m.WithCtx(ctx, h.PingDB())))

	sugar.Infow("Starting server", "addr", flagAddr)

	err = http.ListenAndServe(flagAddr, r)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		panic(err)
	}
}
