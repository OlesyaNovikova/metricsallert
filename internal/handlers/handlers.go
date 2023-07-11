package handlers

type key int

const (
	KeyBD key = iota
)

type MemDataBase interface {
	UpdateGauge(string, float64)
	UpdateCounter(string, int64)
	GetGauge(name string) (value float64, err error)
	GetCounter(name string) (value int64, err error)
	GetAll() map[string]string
}

type MemRepo struct {
	S MemDataBase
}

var memBase MemRepo

func NewMemRepo(Mem MemDataBase) {
	memBase = MemRepo{
		S: Mem,
	}
}
