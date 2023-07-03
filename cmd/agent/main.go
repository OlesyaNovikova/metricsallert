package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	c "github.com/OlesyaNovikova/metricsallert.git/internal/compress"
	j "github.com/OlesyaNovikova/metricsallert.git/internal/json"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
)

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
	Mem.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
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

	return nil
}

func sendJSON(adr string, mem j.Metrics) error {
	b, err := json.Marshal(mem)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, err := c.CompressGzip(b)
	if err != nil {
		fmt.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", adr, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return err
}

func sendMemsJSON(mem s.MemStorage) error {
	var err error
	str := fmt.Sprintf("http://%s/update/", flagAddr)

	for name, val := range mem.MemGauge {
		value := float64(val)
		memJSON := j.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		err = sendJSON(str, memJSON)
		if err != nil {
			fmt.Println(err)
		}
	}
	for name, val := range mem.MemCounter {
		value := int64(val)
		memJSON := j.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}
		err = sendJSON(str, memJSON)
		if err != nil {
			fmt.Println(err)
		}
	}
	return err
}

func main() {
	parseFlags()
	MemBase := s.NewStorage()
	var err error
	tickerP := time.NewTicker(pollInterval)
	defer tickerP.Stop()
	tickerS := time.NewTicker(reportInterval)
	defer tickerS.Stop()
	for {
		select {
		case <-tickerP.C:
			err = collectMems(&MemBase)
			if err != nil {
				fmt.Println(err)
			}
		case <-tickerS.C:
			err = sendMemsJSON(MemBase)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
