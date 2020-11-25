package randutils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandSized(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Rand() string {
	return RandSized(16)
}

func RandList(length int) (randList []string) {
	randSet := make(map[string]bool)
	for len(randList) != length {
		newRand := Rand()
		if _, exists := randSet[newRand]; !exists {
			randSet[newRand] = true
			randList = append(randList, newRand)
		}
	}
	return
}
