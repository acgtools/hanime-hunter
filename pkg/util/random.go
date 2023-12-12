package util

import (
	"math/rand"
	"time"
)

func RandomInt63n(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Int63n(max-min+1) + min
}
