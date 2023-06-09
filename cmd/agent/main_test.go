package main

import (
	"testing"

	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectMems(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  error
		wantLenG int
		wantLenC int
	}{
		{
			name:     "positive case",
			wantLenG: 27,
			wantLenC: 1,
			wantErr:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := s.NewStorage()
			require.Equal(t, test.wantErr, collectMems(&mem))
			assert.Equal(t, test.wantLenG, len(mem.MemGauge))
			assert.Equal(t, test.wantLenC, len(mem.MemCounter))
		})
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name    string
		URL     string
		wantErr error
	}{
		{
			name:    "positive case",
			URL:     "http://localhost:8080/update/",
			wantErr: nil,
		},
		{
			name:    "negative case",
			URL:     "/update/",
			wantErr: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.wantErr, send(test.URL))
		})
	}
}
