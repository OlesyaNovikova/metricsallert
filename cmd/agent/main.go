package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
	g "github.com/OlesyaNovikova/metricsallert.git/internal/utils/gzip"
)

func collectMems(ctx context.Context, Mem *s.MemStorage) error {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	Mem.UpdateCounter(ctx, "PollCount", 1)
	Mem.UpdateGauge(ctx, "RandomValue", rand.Float64())

	Mem.UpdateGauge(ctx, "Alloc", float64(rtm.Alloc))
	Mem.UpdateGauge(ctx, "BuckHashSys", float64(rtm.BuckHashSys))
	Mem.UpdateGauge(ctx, "Frees", float64(rtm.Frees))
	Mem.UpdateGauge(ctx, "GCCPUFraction", float64(rtm.GCCPUFraction))
	Mem.UpdateGauge(ctx, "GCSys", float64(rtm.GCSys))
	Mem.UpdateGauge(ctx, "HeapAlloc", float64(rtm.HeapAlloc))
	Mem.UpdateGauge(ctx, "HeapIdle", float64(rtm.HeapIdle))
	Mem.UpdateGauge(ctx, "HeapInuse", float64(rtm.HeapInuse))
	Mem.UpdateGauge(ctx, "HeapObjects", float64(rtm.HeapObjects))
	Mem.UpdateGauge(ctx, "HeapReleased", float64(rtm.HeapReleased))
	Mem.UpdateGauge(ctx, "HeapSys", float64(rtm.HeapSys))
	Mem.UpdateGauge(ctx, "LastGC", float64(rtm.LastGC))
	Mem.UpdateGauge(ctx, "Lookups", float64(rtm.Lookups))
	Mem.UpdateGauge(ctx, "MCacheInuse", float64(rtm.MCacheInuse))
	Mem.UpdateGauge(ctx, "MCacheSys", float64(rtm.MCacheSys))
	Mem.UpdateGauge(ctx, "MSpanInuse", float64(rtm.MSpanInuse))
	Mem.UpdateGauge(ctx, "MSpanSys", float64(rtm.MSpanSys))
	Mem.UpdateGauge(ctx, "Mallocs", float64(rtm.Mallocs))
	Mem.UpdateGauge(ctx, "NextGC", float64(rtm.NextGC))
	Mem.UpdateGauge(ctx, "NumForcedGC", float64(rtm.NumForcedGC))
	Mem.UpdateGauge(ctx, "NumGC", float64(rtm.NumGC))
	Mem.UpdateGauge(ctx, "OtherSys", float64(rtm.OtherSys))
	Mem.UpdateGauge(ctx, "PauseTotalNs", float64(rtm.PauseTotalNs))
	Mem.UpdateGauge(ctx, "StackInuse", float64(rtm.StackInuse))
	Mem.UpdateGauge(ctx, "StackSys", float64(rtm.StackSys))
	Mem.UpdateGauge(ctx, "Sys", float64(rtm.Sys))
	Mem.UpdateGauge(ctx, "TotalAlloc", float64(rtm.TotalAlloc))

	return nil
}

func sendJSON(ctx context.Context, adr string, mem []j.Metrics) error {
	b, err := json.Marshal(mem)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, err := g.CompressGzip(b)
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

	if !bytes.Equal(KEY, nil) {
		h := hmac.New(sha256.New, KEY)
		h.Write(body)
		dst := hex.EncodeToString(h.Sum(nil))
		fmt.Println(dst)
		req.Header.Add("HashSHA256", dst)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	return err
}

func sendMemsJSON(ctx context.Context, mem *s.MemStorage) error {
	var err error
	str := fmt.Sprintf("http://%s/updates/", flagAddr)

	allMems := []j.Metrics{}

	for name, val := range mem.MemGauge {
		value := float64(val)
		memJSON := j.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		allMems = append(allMems, memJSON)
	}
	for name, val := range mem.MemCounter {
		value := int64(val)
		memJSON := j.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}
		allMems = append(allMems, memJSON)
	}
	delay := 1
	for i := 0; i < 4; i++ {
		err = sendJSON(ctx, str, allMems)
		if err == nil {
			return err
		}
		time.Sleep(time.Duration(delay) * time.Second)
		delay += 2
	}

	fmt.Println(err)
	return err
}

func main() {
	parseFlags()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	MemBase := s.NewStorage()
	var err error
	tickerP := time.NewTicker(pollInterval)
	defer tickerP.Stop()
	tickerS := time.NewTicker(reportInterval)
	defer tickerS.Stop()
	for {
		select {
		case <-tickerP.C:
			err = collectMems(ctx, MemBase)
			if err != nil {
				fmt.Println(err)
			}
		case <-tickerS.C:
			err = sendMemsJSON(ctx, MemBase)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
