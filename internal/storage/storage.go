package storage

import (
	"fmt"
	"strconv"
	"sync"
)

type gauge float64
type counter int64

type MemStorage struct {
	MemGauge   map[string]gauge
	MemCounter map[string]counter
	filePath   string
	saveInFile bool
	mut        sync.Mutex
}

func NewStorage() *MemStorage {
	mem := MemStorage{
		MemGauge:   make(map[string]gauge),
		MemCounter: make(map[string]counter),
		filePath:   "",
		saveInFile: false,
	}
	return &mem
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.MemGauge[name] = gauge(value)
	if m.saveInFile {
		m.mut.Lock()
		m.writeFileStorage()
		m.mut.Unlock()
	}
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.MemCounter[name] += counter(value)
	if m.saveInFile {
		m.mut.Lock()
		m.writeFileStorage()
		m.mut.Unlock()
	}
}

func (m *MemStorage) GetGauge(name string) (value float64, err error) {
	if val, ok := m.MemGauge[name]; ok {
		return float64(val), nil
	}
	err = fmt.Errorf("metric %v not found", name)
	return 0, err
}

func (m *MemStorage) GetCounter(name string) (value int64, err error) {
	if val, ok := m.MemCounter[name]; ok {
		return int64(val), nil
	}
	err = fmt.Errorf("metric %v not found", name)
	return 0, err
}

func (m *MemStorage) getString(name, memtype string) (value string, err error) {

	err = nil
	switch memtype {
	case "gauge":
		if val, ok := m.MemGauge[name]; ok {
			value = strconv.FormatFloat(float64(val), 'f', -1, 64)
			return
		}
		return "", nil

	case "counter":
		if val, ok := m.MemCounter[name]; ok {
			value = strconv.FormatInt(int64(val), 10)
			return
		}
		return "", nil
	}
	err = fmt.Errorf("type %v not found", memtype)
	return "", err
}

func (m *MemStorage) GetAll() map[string]string {

	allMems := make(map[string]string)

	for name := range m.MemGauge {
		allMems[name], _ = m.getString(name, "gauge")
	}
	for name := range m.MemCounter {
		allMems[name], _ = m.getString(name, "counter")
	}
	return allMems
}

func (m *MemStorage) Delete(name, memtype string) {

	switch memtype {
	case "gauge":
		delete(m.MemGauge, name)
	case "counter":
		delete(m.MemCounter, name)
	}
	if m.saveInFile {
		m.mut.Lock()
		m.writeFileStorage()
		m.mut.Unlock()
	}
}
