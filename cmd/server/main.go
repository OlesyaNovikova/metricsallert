package main

import (
	"net/http"

	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
	chi "github.com/go-chi/chi/v5"
)

type MemInterface interface {
	InitStorage()
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetString(name, memtype string) (value string, err error)
	GetAll() map[string]string
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
	r.Post("/update/{memtype}/{name}/{value}", updateMem)
	r.Get("/value/{memtype}/{name}", getMem)
	r.Get("/", getAllMems)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
