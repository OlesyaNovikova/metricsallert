package main

import (
	"flag"
)

var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
