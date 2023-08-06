package handlers

import (
	"bytes"
	"context"
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

		ctx := req.Context()

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
		err = updateJSON(ctx, mem)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(mem)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(resp)
	}
	return http.HandlerFunc(fn)
}

func updateJSON(ctx context.Context, mem j.Metrics) error {
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
		err = memBase.s.UpdateGauge(ctx, mem.ID, *mem.Value)
		if err != nil {
			return err
		}
	case "counter":
		if mem.Delta == nil {
			fmt.Println("BadRequest-value")
			return fmt.Errorf("BadRequest-value")
		}
		del, err := memBase.s.UpdateCounter(ctx, mem.ID, *mem.Delta)
		if err != nil {
			return err
		}
		*mem.Delta = del
	default:
		fmt.Println("BadRequest-type")
		return fmt.Errorf("BadRequest-type")
	}
	return err
}
