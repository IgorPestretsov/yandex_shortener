package storage

import "fmt"

type Storage struct {
	Storage map[string]string
}

func New() *Storage {
	s := Storage{Storage: make(map[string]string)}
	return &s
}

func (s *Storage) LoadLinksPair(key string) string {
	fmt.Println(*s)
	FullLink := s.Storage[key]
	return FullLink
}

func (s *Storage) SaveLinksPair(key string, link string) {
	s.Storage[link] = key
}
