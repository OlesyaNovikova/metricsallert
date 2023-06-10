package main

import (
	"net/http"

	h "github.com/OlesyaNovikova/metricsallert/internal/server/handlers"
	s "github.com/OlesyaNovikova/metricsallert/internal/storage"
	chi "github.com/go-chi/chi/v5"
)

type MemInterface interface {
	InitStorage()
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetString(name, memtype string) (value string, err error)
}

type MemRepo struct {
	S MemInterface
}

var MemBase MemRepo

func init() {
	MemBase = MemRepo{
		S: &s.MemStorage{},
	}
	MemBase.S.InitStorage()
}

func main() {
	r := chi.NewRouter()
	r.Post("/update/{memtype}/{name}/{value}", h.updateMem)
	r.Get("/value/{memtype}/{name}", h.getMem)
	r.Get("/", h.getAllMems)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
