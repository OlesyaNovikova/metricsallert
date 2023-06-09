package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
)

func UpdateMemJSON() http.HandlerFunc {
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
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// десериализуем JSON в Metrics
		if err = json.Unmarshal(inBuf.Bytes(), &mem); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateJSON(mem)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
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

func updateJSON(mem j.Metrics) error {
	var err error

	if mem.ID == "" {
		fmt.Println("BadRequest-name")
		return fmt.Errorf("BadRequest-name")
	}

	switch mem.MType {
	case "gauge":
		if mem.Value == nil {
			fmt.Println("BadRequest-value")
			return fmt.Errorf("BadRequest-value")
		}
		memBase.S.UpdateGauge(mem.ID, *mem.Value)
		*mem.Value, err = memBase.S.GetGauge(mem.ID)
	case "counter":
		if mem.Delta == nil {
			fmt.Println("BadRequest-value")
			return fmt.Errorf("BadRequest-value")
		}
		memBase.S.UpdateCounter(mem.ID, *mem.Delta)
		*mem.Delta, err = memBase.S.GetCounter(mem.ID)
	default:
		fmt.Println("BadRequest-type")
		return fmt.Errorf("BadRequest-type")
	}
	return err
}
