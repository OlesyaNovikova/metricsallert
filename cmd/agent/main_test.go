package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	s "github.com/OlesyaNovikova/metricsallert.git/internal/storage"
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
			wantLenG: 28,
			wantLenC: 1,
			wantErr:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := s.NewStorage()
			ctx := context.Background()
			require.Equal(t, test.wantErr, collectMems(ctx, mem))
			assert.Equal(t, test.wantLenG, len(mem.MemGauge))
			assert.Equal(t, test.wantLenC, len(mem.MemCounter))
		})
	}
}
