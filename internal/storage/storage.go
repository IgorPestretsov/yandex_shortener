package storage

import "fmt"

type Storage struct {
	Storage map[string]string
}

func (s *Storage) Init() {
	s.Storage = make(map[string]string)
}

func (s *Storage) LoadLinksPair(key string) string {
	FullLink := s.Storage[key]
	return FullLink
}

func (s *Storage) SaveLinksPair(key string, link string) {
	s.Storage[link] = key
	fmt.Println()
}
