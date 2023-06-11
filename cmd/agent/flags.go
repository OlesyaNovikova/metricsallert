package main

import (
	"flag"
	"time"
)

var flagAddr string

var pollInterval time.Duration
var reportInterval time.Duration

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	p := flag.Int64("p", 2, "metric collection interval")
	r := flag.Int64("r", 10, "metrics sending interval")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	pollInterval = time.Duration(*p) * time.Second
	reportInterval = time.Duration(*r) * time.Second
}
