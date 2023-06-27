package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

type MemDataBase interface {
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetGauge(name string) (value float64, err error)
	GetCounter(name string) (value int64, err error)
	GetString(name, memtype string) (value string, err error)
	GetAll() map[string]string
}

type MemRepo struct {
	S MemDataBase
}

var memBase MemRepo

func NewMemRepo(Mem MemDataBase) {
	memBase = MemRepo{
		S: Mem,
	}
}

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

func GetMem() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			fmt.Print("Only GET requests are allowed!\n")
			http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		memtype := chi.URLParam(req, "memtype")
		name := chi.URLParam(req, "name")

		strValue, err := memBase.S.GetString(name, memtype)

		if err != nil {
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
