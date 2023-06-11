package main

import (
	"flag"
	"os"
)

var flagAddr string

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		flagAddr = envAddr
	}
}
