package random

import (
	"math/rand"
	"time"
)

var (
	defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[defaultRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomNumber(min, max int) int  {
	return defaultRand.Intn(max-min) + min
}