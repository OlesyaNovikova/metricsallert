package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdMem(t *testing.T) {

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		name         string
		method       string
		target       string
		expectedCode int
		expectedBody string
	}{
		{name: "Метод GET", method: http.MethodGet, target: "/update/gauge/HeapSys/32", expectedCode: http.StatusMethodNotAllowed, expectedBody: "Only POST requests are allowed!\n"},
		{name: "Корректные данные(gauge)", method: http.MethodPost, target: "/update/gauge/HeapSys/32", expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Корректные данные(counter)", method: http.MethodPost, target: "/update/counter/HeapSys/32", expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Не верный тип", method: http.MethodPost, target: "/update/try/HeapSys/32", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Не верное значение", method: http.MethodPost, target: "/update/gauge/HeapSys/try", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Не задано имя метрики", method: http.MethodPost, target: "/update/gauge/", expectedCode: http.StatusNotFound, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(tc.method, tc.target, nil)
			w := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/update/{memtype}/{name}/{meaning}", updMem)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
