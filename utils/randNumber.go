package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func RandNumber(max int) string {
	rand.Seed(time.Now().Unix())

	//Generate a random array of length n
	list := rand.Perm(10)
	var intS string

	for _, r := range list {
		intS += strconv.Itoa(r)

		if len(intS) == max {
			break
		}
	}

	return intS

}

func RandPass(max int) string {
	var res string
	var letArr = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < max; i++ {
		res += letArr[r1.Intn(len(letArr))]
	}

	return res
}
