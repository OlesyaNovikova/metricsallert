package main

import (
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

func updateMem(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run updMem:\n")
	if req.Method != http.MethodPost {
		fmt.Print("Only POST requests are allowed!\n")
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	vars := chi.Vars(req)
	memtype := vars["memtype"]
	name := vars["name"]
	value := vars["value"]

	switch memtype {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Println("BadRequest-meaning")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		MemBase.S.UpdateGauge(name, val)
		fmt.Println(MemBase)
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
		fmt.Println(MemBase)
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

	vars := chi.Vars(req)
	memtype := vars["memtype"]
	name := vars["name"]

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

func getAllMems(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run getAllMems:\n")
	if req.Method != http.MethodGet {
		fmt.Print("Only GET requests are allowed!\n")
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var str string
	///////

	///////
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(str))
}
