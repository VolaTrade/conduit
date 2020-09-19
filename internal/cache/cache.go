package cache

import (
	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type (
	Cache interface {
	}

	CandlesCache struct {
		value string
	}
)

func New() *CandlesCache {

	return &CandlesCache{}

}
