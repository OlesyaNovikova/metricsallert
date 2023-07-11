package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

func UpdateMem() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			fmt.Print("Only POST requests are allowed!\n")
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		memtype := chi.URLParam(req, "memtype")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")

		switch memtype {
		case "gauge":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("BadRequest-meaning")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			memBase.S.UpdateGauge(name, val)
			res.WriteHeader(http.StatusOK)
			return

		case "counter":
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				fmt.Println("BadRequest-meaning")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			memBase.S.UpdateCounter(name, val)
			res.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("BadRequest-type")
		res.WriteHeader(http.StatusBadRequest)
	}
	return http.HandlerFunc(fn)
}
