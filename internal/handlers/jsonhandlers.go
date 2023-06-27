package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/json"
)

func UpdateMemJson() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			fmt.Print("Only POST requests are allowed!\n")
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
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

		switch mem.MType {
		case "gauge":
			memBase.S.UpdateGauge(mem.ID, *mem.Value)
			*mem.Value, err = memBase.S.GetGauge(mem.ID)
		case "counter":
			memBase.S.UpdateCounter(mem.ID, *mem.Delta)
			*mem.Delta, err = memBase.S.GetCounter(mem.ID)
		default:
			fmt.Println("BadRequest-type")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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

func GetMemJson() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			fmt.Print("Only GET requests are allowed!\n")
			http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

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

		switch mem.MType {
		case "gauge":
			*mem.Value, err = memBase.S.GetGauge(mem.ID)
		case "counter":
			*mem.Delta, err = memBase.S.GetCounter(mem.ID)
		default:
			fmt.Println("BadRequest-type")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			fmt.Println("Metric not found")
			res.WriteHeader(http.StatusNotFound)
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
