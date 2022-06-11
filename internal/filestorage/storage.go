package filestorage

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"log"
)

type NestedMap map[string]map[string]string

type Storage struct {
	Storage NestedMap
	r       *reader
	w       *writer
	cleaner Cleaner
}

func (s *Storage) GetChannelForDelete() chan map[string]string {
	return s.cleaner.UserDeleteRequests
}

func New(filepath string, q chan bool) *Storage {
	data := make(NestedMap)
	s := Storage{Storage: data}
	s.cleaner = NewCleaner(&s)
	s.cleaner.Run(q)
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

		if fullLink, ok := value[key]; ok {
			return fullLink
		}
	}
	return ""
}

func (s *Storage) SaveLinksPair(uid string, link string, key string) (string, error) {
	if _, ok := s.Storage[uid]; !ok {
		s.Storage[uid] = make(map[string]string)
	}
	s.Storage[uid][key] = link
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
func (s *Storage) DeleteRecord(toDelete RecordToDelete) {
	s.Storage[toDelete.userID][toDelete.urlID] = storage.GoneValue
}
