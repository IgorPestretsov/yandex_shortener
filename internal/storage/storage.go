package storage

import (
	"fmt"
)

var (
	KeyChannel        chan string
	FullLinkChannel   chan string
	LinksPairsChannel chan [2]string
	Storage           map[string]string
)

func LoadLinksPair() {
	for {
		key := <-KeyChannel
		FullLink := Storage[key]
		FullLinkChannel <- FullLink

	}
}

func SaveLinksPair() {
	for {
		LinksPair := <-LinksPairsChannel
		Storage[LinksPair[1]] = LinksPair[0]
		fmt.Println(Storage)
	}
}
