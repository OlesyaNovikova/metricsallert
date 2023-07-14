package main

import (
	"flag"
	"os"
	"strconv"
	"time"
)

var (
	flagAddr        string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	DBAddr          string
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	s := flag.Int64("i", 300, "time interval to save to disk")
	flag.StringVar(&FileStoragePath, "f", "./tmp/metrics-db.json", "file where the current values are saved")
	flag.BoolVar(&Restore, "r", true, "load previously saved values from the specified file")
	flag.StringVar(&DBAddr, "d", "", "data base DSN")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		flagAddr = envAddr
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if envS, err := strconv.Atoi(envStoreInterval); err == nil {
			*s = int64(envS)
		}
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		FileStoragePath = envStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if envRestore == "true" || envRestore == "TRUE" {
			Restore = true
		} else if envRestore == "false" || envRestore == "FALSE" {
			Restore = false
		}
	}
	StoreInterval = time.Duration(*s) * time.Second

	if envDBAddr := os.Getenv("DATABASE_DSN"); envDBAddr != "" {
		DBAddr = envDBAddr
	}
}
