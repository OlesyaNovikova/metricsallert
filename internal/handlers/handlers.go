package handlers

import (
	"context"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
)

type MemRepository interface {
	UpdateGauge(context.Context, string, float64) error
	UpdateCounter(context.Context, string, int64) (int64, error)
	GetGauge(ctx context.Context, name string) (value float64, err error)
	GetCounter(ctx context.Context, name string) (value int64, err error)
	GetAll(ctx context.Context) (map[string]string, error)
	Ping(ctx context.Context) error
	Updates(ctx context.Context, mems []j.Metrics) error
}

type MemRepo struct {
	s MemRepository
}

var memBase MemRepo

func NewMemRepo(Mem MemRepository) {
	memBase = MemRepo{
		s: Mem,
	}
}
