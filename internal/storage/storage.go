package storage

import (
	"fmt"
	"strconv"
)

type gauge float64
type counter int64

type MemStorage struct {
	MemGauge   map[string]gauge
	MemCounter map[string]counter
}

func NewStorage() (m MemStorage) {
	return MemStorage{
		MemGauge:   make(map[string]gauge),
		MemCounter: make(map[string]counter),
	}
}

func (m *MemStorage) InitStorage() {
	m.MemGauge = make(map[string]gauge)
	m.MemCounter = make(map[string]counter)
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.MemGauge[name] = gauge(value)
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.MemCounter[name] += counter(value)
}

func (m *MemStorage) GetString(name, memtype string) (value string, err error) {

	err = nil
	switch memtype {
	case "gauge":
		if val, ok := m.MemGauge[name]; ok {
			value = strconv.FormatFloat(float64(val), 'f', 5, 64)
			return
		}
		return "", nil

	case "counter":
		if val, ok := m.MemGauge[name]; ok {
			value = strconv.FormatInt(int64(val), 10)
			return
		}
		return "", nil
	}
	err = fmt.Errorf("type %v not found", memtype)
	return "", err
}

func (m *MemStorage) Delete(name, memtype string) {

	switch memtype {
	case "gauge":
		delete(m.MemGauge, name)
	case "counter":
		delete(m.MemCounter, name)
	}
}

/*
type MemInterface interface {
	InitStorage()
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetString(name, memtype string) (value string, err error)
	Delete(name, memtype string)
}
*/
