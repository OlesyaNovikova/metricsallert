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
	flag.StringVar(&flagAddr, "a", ":8080", "address and port to run server")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "metric collection interval")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "metrics sending interval")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
