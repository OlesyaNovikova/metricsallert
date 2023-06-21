package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"

	h "github.com/OlesyaNovikova/metricsallert.git/internal/handlers"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

func main() {

	parseFlags()

	mem := s.NewStorage()
	h.NewMemRepo(&mem)

	r := chi.NewRouter()
	r.Post("/update/{memtype}/{name}/{value}", h.UpdateMem)
	r.Get("/value/{memtype}/{name}", h.GetMem)
	r.Get("/", h.GetAllMems)

	err := http.ListenAndServe(flagAddr, r)
	if err != nil {
		panic(err)
	}
}
