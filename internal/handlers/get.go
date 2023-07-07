package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

func GetMem() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			fmt.Print("Only GET requests are allowed!\n")
			http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		memtype := chi.URLParam(req, "memtype")
		name := chi.URLParam(req, "name")

		var strValue string
		switch memtype {
		case "gauge":
			valF, err := memBase.S.GetGauge(name)
			if err != nil {
				fmt.Println("Metric not found")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			strValue = strconv.FormatFloat(float64(valF), 'f', -1, 64)
		case "counter":
			valI, err := memBase.S.GetCounter(name)
			if err != nil {
				fmt.Println("Metric not found")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			strValue = strconv.FormatInt(int64(valI), 10)
		default:
			fmt.Println("BadRequest-type")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if strValue == "" {
			fmt.Println("Metric not found")
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strValue))
	}
	return http.HandlerFunc(fn)
}
