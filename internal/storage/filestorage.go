package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
)

func NewFileStorage(ctx context.Context, path string, restore bool, interval time.Duration) (*MemStorage, error) {
	var err error
	err = nil
	mem := NewStorage()

	if _, err = os.Stat(path); err == nil {
		mem.filePath = path
		if restore {
			err = mem.readFileStorage(ctx)
		}
	} else {
		file, err := os.Create(path)
		if err != nil {
			return mem, err
		}
		file.Close()
		mem.filePath = path
	}
	if interval == 0 {
		mem.saveInFile = true
	} else {
		go mem.fileStorageRoutine(ctx, interval)
	}

	return mem, err
}

func (m *MemStorage) fileStorageRoutine(ctx context.Context, interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("контекст отменен")
			return
		case <-ticker.C:
			m.mut.Lock()
			err := m.writeFileStorage()
			m.mut.Unlock()
			if err != nil {
				fmt.Printf("file write error(%v)", err)
				return
			}
		}
	}
}

func (m *MemStorage) FileStorageExit() {
	err := m.writeFileStorage()
	if err != nil {
		fmt.Printf("file write error(%v)", err)
		return
	}
}

func (m *MemStorage) readFileStorage(ctx context.Context) error {
	file, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var mem j.Metrics
		err := json.Unmarshal(data, &mem)
		if err != nil {
			return err
		}
		err = m.updateJSON(ctx, mem)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) writeFileStorage() error {
	file, err := os.OpenFile(m.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	mems := m.getAllForJSON()

	for _, mem := range mems {
		data, err := json.Marshal(&mem)
		if err != nil {
			return err
		}
		// добавляем перенос строки
		data = append(data, '\n')
		_, err = file.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) getAllForJSON() []j.Metrics {
	mems := make([]j.Metrics, 0)
	for name, val := range m.MemGauge {
		value := float64(val)
		mem := j.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		mems = append(mems, mem)
	}
	for name, del := range m.MemCounter {
		delta := int64(del)
		mem := j.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		mems = append(mems, mem)
	}
	return mems
}

func (m *MemStorage) updateJSON(ctx context.Context, mem j.Metrics) error {

	if mem.ID == "" {
		return fmt.Errorf("no name")
	}

	switch mem.MType {
	case "gauge":
		if mem.Value == nil {
			return fmt.Errorf("no value")
		}
		m.UpdateGauge(ctx, mem.ID, *mem.Value)

	case "counter":
		if mem.Delta == nil {
			return fmt.Errorf("no delta")
		}
		m.UpdateCounter(ctx, mem.ID, *mem.Delta)
	default:
		return fmt.Errorf("bad type")
	}
	return nil
}
