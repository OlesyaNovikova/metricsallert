package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type mem struct {
	gauge   float64
	counter int64
}

type MemStorage struct {
	Mem map[string]mem
}

var MemBase MemStorage

func updMem(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run updMem:\n")

	if req.Method != http.MethodPost {
		fmt.Print("Only POST requests are allowed!\n")
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(req)
	memtype := vars["memtype"]
	name := vars["name"]
	smeaning := vars["meaning"]

	if memtype == "gauge" {
		meaning, err := strconv.ParseFloat(smeaning, 64)
		if err != nil {
			fmt.Print("BadRequest-meaning\n")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		addmem := MemBase.Mem[name]
		addmem.gauge = meaning
		MemBase.Mem[name] = addmem
		fmt.Print(MemBase)
		res.WriteHeader(http.StatusOK)
		return

	} else if memtype == "counter" {
		meaning, err := strconv.ParseInt(smeaning, 10, 64)
		if err != nil {
			fmt.Print("BadRequest-meaning\n")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		addmem := MemBase.Mem[name]
		addmem.counter += meaning
		MemBase.Mem[name] = addmem
		fmt.Print(MemBase)
		res.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Print("BadRequest-type\n")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	MemBase.Mem = make(map[string]mem)

	r := mux.NewRouter()
	r.HandleFunc("/update/{memtype}/{name}/{meaning}", updMem)
	r.HandleFunc("/update/{memtype}/{name}/{meaning}/", updMem)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
