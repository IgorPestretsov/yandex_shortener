package storage

import (
	"fmt"
	"log"
)

type Storage struct {
	Storage  map[string]string
	consumer *consumer
	producer *producer
}

func New(filepath string) *Storage {
	data := make(map[string]string)
	s := Storage{Storage: data}
	if filepath != "" {
		p, err := NewProducer(filepath)
		if err != nil {
			log.Fatal(err)
		}
		s.producer = p

		c, err := NewConsumer(filepath)
		if err != nil {
			log.Fatal(err)
		}
		s.consumer = c

		s.loadStorageFromFile()
	}

	return &s
}

func (s *Storage) loadStorageFromFile() {
	var err error
	s.Storage, err = s.consumer.ReadData()
	if err != nil {
		log.Println("Cannot read storage file. ")
	}
}

func (s *Storage) LoadLinksPair(key string) string {
	FullLink := s.Storage[key]
	return FullLink
}

func (s *Storage) SaveLinksPair(key string, link string) {
	s.Storage[link] = key
	if s.producer != nil {
		toFile := map[string]string{link: key}
		err := s.producer.WriteEvent(toFile)

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(s)

}

func (s *Storage) Close() {
	s.producer.Close()
	s.consumer.Close()
}
