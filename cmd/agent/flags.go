package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	flagAddr       string
	pollInterval   time.Duration
	reportInterval time.Duration
	KEY            []byte
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	p := flag.Int64("p", 2, "metric collection interval")
	r := flag.Int64("r", 10, "metrics sending interval")
	k := flag.String("k", "", "key")
	// парсим переданны аргументы в зарегистрированные переменные
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		flagAddr = envAddr
	}

	if envR := os.Getenv("POLL_INTERVAL"); envR != "" {
		if rEnv, err := strconv.Atoi(envR); err == nil {
			*r = int64(rEnv)
		}
	}

	if envP := os.Getenv("REPORT_INTERVAL"); envP != "" {
		if pEnv, err := strconv.Atoi(envP); err == nil {
			*p = int64(pEnv)
		}
	}

	pollInterval = time.Duration(*p) * time.Second
	reportInterval = time.Duration(*r) * time.Second

	if envKey := os.Getenv("KEY"); envKey != "" {
		fmt.Printf("Переменная окружения %v \n", envKey)
		*k = envKey
	}
	if *k != "" {
		var err error
		fmt.Println(*k)
		KEY, err = hex.DecodeString(*k)
		if err != nil {
			panic(err)
		}
	}
}
