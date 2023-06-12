package main

import (
	"net/http"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
	chi "github.com/go-chi/chi/v5"
)

func main() {

	parseFlags()

	h.NewMemRepo(&s.MemStorage{})

	r := chi.NewRouter()
	r.Post("/update/{memtype}/{name}/{value}", h.UpdateMem)
	r.Get("/value/{memtype}/{name}", h.GetMem)
	r.Get("/", h.GetAllMems)

	err := http.ListenAndServe(flagAddr, r)
	if err != nil {
		panic(err)
	}
}
