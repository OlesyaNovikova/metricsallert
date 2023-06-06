package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

const servAdr string = "http://localhost:8080/update/"
const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

type gauge float64
type counter int64

type MemStorage struct {
	MemGauge   map[string]gauge
	MemCounter map[string]counter
}

func collectMems(Mem *MemStorage) error {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	Mem.MemGauge["Alloc"] = gauge(rtm.Alloc)
	Mem.MemGauge["BuckHashSys"] = gauge(rtm.BuckHashSys)
	Mem.MemGauge["Frees"] = gauge(rtm.Frees)
	Mem.MemGauge["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
	Mem.MemGauge["GCSys"] = gauge(rtm.GCSys)
	Mem.MemGauge["HeapAlloc"] = gauge(rtm.HeapAlloc)
	Mem.MemGauge["HeapIdle"] = gauge(rtm.HeapIdle)
	Mem.MemGauge["HeapInuse"] = gauge(rtm.HeapInuse)
	Mem.MemGauge["HeapObjects"] = gauge(rtm.HeapObjects)
	Mem.MemGauge["HeapReleased"] = gauge(rtm.HeapReleased)
	Mem.MemGauge["HeapSys"] = gauge(rtm.HeapSys)
	Mem.MemGauge["LastGC"] = gauge(rtm.LastGC)
	Mem.MemGauge["Lookups"] = gauge(rtm.Lookups)
	Mem.MemGauge["MCacheInuse"] = gauge(rtm.MCacheInuse)
	Mem.MemGauge["MCacheSys"] = gauge(rtm.MCacheSys)
	Mem.MemGauge["MSpanInuse"] = gauge(rtm.MSpanInuse)
	Mem.MemGauge["Mallocs"] = gauge(rtm.Mallocs)
	Mem.MemGauge["NextGC"] = gauge(rtm.NextGC)
	Mem.MemGauge["NumForcedGC"] = gauge(rtm.NumForcedGC)
	Mem.MemGauge["NumGC"] = gauge(rtm.NumGC)
	Mem.MemGauge["OtherSys"] = gauge(rtm.OtherSys)
	Mem.MemGauge["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	Mem.MemGauge["StackInuse"] = gauge(rtm.StackInuse)
	Mem.MemGauge["StackSys"] = gauge(rtm.StackSys)
	Mem.MemGauge["Sys"] = gauge(rtm.Sys)
	Mem.MemGauge["TotalAlloc"] = gauge(rtm.TotalAlloc)
	//
	Mem.MemCounter["PollCount"] += 1
	Mem.MemGauge["RandomValue"] = gauge(rand.Float64())

	// Just encode to json and print
	b, _ := json.Marshal(Mem)
	fmt.Println(string(b))

	return nil
}

func send(adr string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", adr, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return err
}

func sendMems(mem MemStorage) error {
	var str string
	var err error
	for name, val := range mem.MemGauge {
		value := strconv.FormatFloat(float64(val), 'f', 5, 64)
		str = fmt.Sprintf("%sgauge/%s/%s", servAdr, name, value)
		fmt.Println(str)
		err = send(str)
		if err != nil {
			fmt.Println(err)
			//return err
		}
	}
	for name, val := range mem.MemCounter {
		value := strconv.FormatInt(int64(val), 10)
		str = fmt.Sprintf("%scounter/%s/%s", servAdr, name, value)
		fmt.Print(str)
		err = send(str)
		if err != nil {
			fmt.Println(err)
			//return err
		}
	}
	return err
}

func main() {
	var MemBase MemStorage
	MemBase.MemGauge = make(map[string]gauge)
	MemBase.MemCounter = make(map[string]counter)
	var err error
	var time_t time.Duration
	for {
		err = collectMems(&MemBase)
		if err != nil {
			fmt.Println(err)
			//return err
		}
		time.Sleep(pollInterval)
		time_t += pollInterval
		if time_t >= reportInterval {
			time_t = 0
			err = sendMems(MemBase)
			if err != nil {
				fmt.Println(err)
				//return err
			}
		}
	}
}
