package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/mem"

	js "github.com/OlesyaNovikova/metricsallert.git/internal/models"
	st "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
	gz "github.com/OlesyaNovikova/metricsallert.git/internal/utils/gzip"
)

type gauge struct {
	name  string
	value float64
}

type counter struct {
	name  string
	delta int64
}

func main() {
	parseFlags()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan counter, 1)
	g := make(chan gauge, 30)

	var snc sync.WaitGroup
	snc.Add(3)

	go readGopsutil(ctx, pollInterval, g, &snc)
	go readMemStats(ctx, pollInterval, c, g, &snc)
	go collectMems(ctx, reportInterval, c, g, rateLimit, &snc)

	snc.Wait()
}

func readGopsutil(ctx context.Context, pollInt time.Duration, g chan gauge, snc *sync.WaitGroup) {
	defer snc.Done()
	ticker := time.NewTicker(pollInt)
	defer ticker.Stop()
	//proc := runtime.NumCPU()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return
			}
			g <- gauge{name: "TotalMemory", value: float64(memInfo.Total)}
			g <- gauge{name: "FreeMemory", value: float64(memInfo.Free)}

			//g <- gauge{name: , value: }
		}
	}
}

func readMemStats(ctx context.Context, pollInt time.Duration, c chan counter, g chan gauge, snc *sync.WaitGroup) {
	defer snc.Done()
	ticker := time.NewTicker(pollInt)
	defer ticker.Stop()
	var rtm runtime.MemStats

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runtime.ReadMemStats(&rtm)
			c <- counter{name: "PollCount", delta: 1}
			g <- gauge{name: "RandomValue", value: rand.Float64()}
			g <- gauge{name: "Alloc", value: float64(rtm.Alloc)}
			g <- gauge{name: "BuckHashSys", value: float64(rtm.BuckHashSys)}
			g <- gauge{name: "Frees", value: float64(rtm.Frees)}
			g <- gauge{name: "GCCPUFraction", value: float64(rtm.GCCPUFraction)}
			g <- gauge{name: "GCSys", value: float64(rtm.GCSys)}
			g <- gauge{name: "HeapAlloc", value: float64(rtm.HeapAlloc)}
			g <- gauge{name: "HeapIdle", value: float64(rtm.HeapIdle)}
			g <- gauge{name: "HeapInuse", value: float64(rtm.HeapInuse)}
			g <- gauge{name: "HeapObjects", value: float64(rtm.HeapObjects)}
			g <- gauge{name: "HeapReleased", value: float64(rtm.HeapReleased)}
			g <- gauge{name: "HeapSys", value: float64(rtm.HeapSys)}
			g <- gauge{name: "LastGC", value: float64(rtm.LastGC)}
			g <- gauge{name: "Lookups", value: float64(rtm.Lookups)}
			g <- gauge{name: "MCacheInuse", value: float64(rtm.MCacheInuse)}
			g <- gauge{name: "MCacheSys", value: float64(rtm.MCacheSys)}
			g <- gauge{name: "MSpanInuse", value: float64(rtm.MSpanInuse)}
			g <- gauge{name: "MSpanSys", value: float64(rtm.MSpanSys)}
			g <- gauge{name: "Mallocs", value: float64(rtm.Mallocs)}
			g <- gauge{name: "NextGC", value: float64(rtm.NextGC)}
			g <- gauge{name: "NumForcedGC", value: float64(rtm.NumForcedGC)}
			g <- gauge{name: "NumGC", value: float64(rtm.NumGC)}
			g <- gauge{name: "OtherSys", value: float64(rtm.OtherSys)}
			g <- gauge{name: "PauseTotalNs", value: float64(rtm.PauseTotalNs)}
			g <- gauge{name: "StackInuse", value: float64(rtm.StackInuse)}
			g <- gauge{name: "StackSys", value: float64(rtm.StackSys)}
			g <- gauge{name: "Sys", value: float64(rtm.Sys)}
			g <- gauge{name: "TotalAlloc", value: float64(rtm.TotalAlloc)}
		}
	}
}

func collectMems(ctx context.Context, repInt time.Duration, c chan counter, g chan gauge, lim int, snc *sync.WaitGroup) {
	defer snc.Done()
	ticker := time.NewTicker(repInt)
	defer ticker.Stop()
	mem := st.NewStorage()
	if lim <= 0 {
		lim = 1
	}
	m := make(chan []js.Metrics, lim)
	for l := 0; l < lim; l++ {
		go sendMemsJSON(ctx, m)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case cou := <-c:
			mem.UpdateCounter(ctx, cou.name, cou.delta)
		case gau := <-g:
			mem.UpdateGauge(ctx, gau.name, gau.value)
		case <-ticker.C:
			allMems := mem.GetAllForJSON()
			m <- allMems
		}
	}
}

func sendMemsJSON(ctx context.Context, m chan []js.Metrics) error {
	var err error
	str := fmt.Sprintf("http://%s/updates/", flagAddr)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("end of context")
			return err
		case allMems := <-m:
			delay := 1
			for i := 0; i < 4; i++ {
				err = sendJSON(ctx, str, allMems)
				if err == nil {
					return err
				}
				time.Sleep(time.Duration(delay) * time.Second)
				delay += 2
			}
		}
	}
}

func sendJSON(ctx context.Context, adr string, mem []js.Metrics) error {
	b, err := json.Marshal(mem)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, err := gz.CompressGzip(b)
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

	if KEY != nil {
		h := hmac.New(sha256.New, KEY)
		h.Write(body)
		dst := h.Sum(nil)
		req.Header.Add("HashSHA256", string(dst))
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
