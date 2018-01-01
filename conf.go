package main

import (
	"encoding/json"
	"log"
	"os"
)

type Conf struct {
	Addr string
}

func GetConf(f string) (c Conf) {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal("Error opening conf file:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("Error decoding conf file:", err)
	}
	return
}
