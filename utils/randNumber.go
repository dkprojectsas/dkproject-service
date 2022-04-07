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
