package utils

import (
	"math/rand"
)

func Shuffle[T comparable](a []T) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}
