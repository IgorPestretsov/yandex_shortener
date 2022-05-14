package storage

import (
	"log"
)

type Storage struct {
	Storage map[string]string
	r       *reader
	w       *writer
}

func New(filepath string) *Storage {
	data := make(map[string]string)
	s := Storage{Storage: data}
	if filepath != "" {
		w, err := NewWriter(filepath)
		if err != nil {
			log.Fatal(err)
		}
		s.w = w

		r, err := NewReader(filepath)
		if err != nil {
			log.Fatal(err)
		}
		s.r = r

		s.loadStorageFromFile()
	}

	return &s
}

func (s *Storage) loadStorageFromFile() {
	var err error
	s.Storage, err = s.r.ReadData()
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
	if s.w != nil {
		toFile := map[string]string{link: key}
		err := s.w.WriteEvent(toFile)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func (s *Storage) Close() {
	s.w.Close()
	s.r.Close()
}
