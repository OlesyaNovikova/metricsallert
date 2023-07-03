package handlers

import (
	"bufio"
	"encoding/json"
	"os"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/json"
)

func ReadFileStorage(path string) error {
	file, err := os.Open(path)
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
		err = updateJSON(mem)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func WriteFileStorage(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	memBase.mut.Lock()
	mems := memBase.S.GetAllForJSON()
	memBase.mut.Unlock()

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
