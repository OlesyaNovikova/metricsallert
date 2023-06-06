package main

import (
	"testing"

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
			var mem MemStorage
			mem.MemGauge = make(map[string]gauge)
			mem.MemCounter = make(map[string]counter)
			require.Equal(t, test.wantErr, collectMems(&mem))
			assert.Equal(t, test.wantLenG, len(mem.MemGauge))
			assert.Equal(t, test.wantLenC, len(mem.MemCounter))
		})
	}
}
