package FileStorage

import (
	"log"
)

type Storage struct {
	Storage map[string]map[string]string
	r       *reader
	w       *writer
}

func New(filepath string) *Storage {
	data := make(map[string]map[string]string)
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
		log.Println("Cannot read FileStorage file. ")
	}
}

func (s *Storage) LoadLinksPair(key string) string {
	for _, value := range s.Storage {
		if FullLink, ok := value[key]; ok {
			return FullLink
		}
	}
	return ""
}

func (s *Storage) SaveLinksPair(uid string, key string, link string) (string, error) {
	if _, ok := s.Storage[uid]; !ok {
		s.Storage[uid] = make(map[string]string)
	}
	s.Storage[uid][link] = key
	if s.w != nil {
		err := s.w.WriteEvent(s.Storage)

		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (s *Storage) GetAllUserURLs(uid string) map[string]string {
	return s.Storage[uid]
}

func (s *Storage) Close() {
	s.w.Close()
	s.r.Close()
}
