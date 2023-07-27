package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
)

func GetMemJSON() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			fmt.Print("Only POST requests are allowed!\n")
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		ctx := req.Context()
		var mem j.Metrics
		var inBuf bytes.Buffer
		// читаем тело запроса
		_, err := inBuf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		// десериализуем JSON в Metrics
		if err = json.Unmarshal(inBuf.Bytes(), &mem); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if mem.ID == "" {
			fmt.Println("BadRequest-name")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var valF float64
		var valI int64

		switch mem.MType {
		case "gauge":
			valF, err = memBase.s.GetGauge(ctx, mem.ID)
			if err == nil {
				mem.Value = &valF
			} else {
				fmt.Println("Metric not found")
				res.WriteHeader(http.StatusNotFound)
				return
			}
		case "counter":
			valI, err = memBase.s.GetCounter(ctx, mem.ID)
			if err == nil {
				mem.Delta = &valI
			} else {
				fmt.Println("Metric not found")
				res.WriteHeader(http.StatusNotFound)
				return
			}
		default:
			fmt.Println("BadRequest-type")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(mem)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
	}
	return http.HandlerFunc(fn)
}
