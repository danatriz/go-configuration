package main

import (
	"configurations/test"
	"log"
)

func main() {
	graylog, err := test.TestGraylog()
	if err != nil {
		panic(err)
	}
	log.Println(graylog)
}
