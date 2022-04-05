package config

import (
	"fmt"
)

func FailOnError(err error, line int, fileName string) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Printf("%v - %s - %s", line, fileName, err.Error())
		fmt.Println(err)
		return
	}
}
