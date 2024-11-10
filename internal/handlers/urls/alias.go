package urls

import (
	"fmt"
	"math"
	"sync/atomic"
)

const (
	seq         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789()@:%_+.~#&="
	seqLen      = uint32(len(seq))
	maxAliasLen = 5
)

var (
	maxCount                 = uint32(math.Pow(float64(seqLen), maxAliasLen))
	ErrMaxAliasCountEexceeds = fmt.Errorf("alias count exceeds the maximum permissible value: %d", maxCount)
)

type aliasService struct {
	count *atomic.Uint32
}

func newAliasService(currAliasCount uint32) (aliasService, error) {
	if currAliasCount >= maxCount {
		return aliasService{}, ErrMaxAliasCountEexceeds
	}
	count := &atomic.Uint32{}
	count.Store(currAliasCount)
	return aliasService{count}, nil
}

func (s *aliasService) nextAlias() (string, error) {
	currAlias := s.count.Add(1) - 1
	if currAlias >= maxCount {
		return "", ErrMaxAliasCountEexceeds
	}

	return getAlias(currAlias), nil
}

func getAlias(i uint32) string {
	if i < seqLen {
		return string(seq[i])
	}
	return getAlias(i/seqLen-1) + string(seq[i%seqLen])
}

func (s *aliasService) loadCount() uint32 {
	return s.count.Load()
}
