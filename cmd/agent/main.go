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

	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

const servAdr string = "http://localhost:8080/update/"

func collectMems(Mem *s.MemStorage) error {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	Mem.UpdateCounter("PollCount", 1)
	Mem.UpdateGauge("RandomValue", rand.Float64())

	Mem.UpdateGauge("Alloc", float64(rtm.Alloc))
	Mem.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	Mem.UpdateGauge("Frees", float64(rtm.Frees))
	Mem.UpdateGauge("GCCPUFraction", float64(rtm.GCCPUFraction))
	Mem.UpdateGauge("GCSys", float64(rtm.GCSys))
	Mem.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	Mem.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	Mem.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	Mem.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	Mem.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	Mem.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	Mem.UpdateGauge("LastGC", float64(rtm.LastGC))
	Mem.UpdateGauge("Lookups", float64(rtm.Lookups))
	Mem.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	Mem.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	Mem.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	Mem.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	Mem.UpdateGauge("NextGC", float64(rtm.NextGC))
	Mem.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	Mem.UpdateGauge("NumGC", float64(rtm.NumGC))
	Mem.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	Mem.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	Mem.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	Mem.UpdateGauge("StackSys", float64(rtm.StackSys))
	Mem.UpdateGauge("Sys", float64(rtm.Sys))
	Mem.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	//
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

func sendMems(mem s.MemStorage) error {
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
	parseFlags()
	MemBase := s.NewStorage()
	var err error
	var timeT time.Duration
	for {
		err = collectMems(&MemBase)
		if err != nil {
			fmt.Println(err)
			//return err
		}
		time.Sleep(pollInterval)
		timeT += pollInterval
		if timeT >= reportInterval {
			timeT = 0
			err = sendMems(MemBase)
			if err != nil {
				fmt.Println(err)
				//return err
			}
		}
	}
}
