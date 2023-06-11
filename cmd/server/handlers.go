package main

import (
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

func updateMem(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run updateMem:\n")
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
		MemBase.S.UpdateGauge(name, val)
		res.WriteHeader(http.StatusOK)
		return

	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Println("BadRequest-meaning")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		MemBase.S.UpdateCounter(name, val)
		res.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("BadRequest-type")
	res.WriteHeader(http.StatusBadRequest)
	return
}

func getMem(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run getMem:\n")
	if req.Method != http.MethodGet {
		fmt.Print("Only GET requests are allowed!\n")
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	memtype := chi.URLParam(req, "memtype")
	name := chi.URLParam(req, "name")

	strValue, err := MemBase.S.GetString(name, memtype)

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
