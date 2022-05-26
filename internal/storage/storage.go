package storage

type Storage interface {
	LoadLinksPair(key string) string
	SaveLinksPair(uid string, key string, link string) (string, error)
	GetAllUserURLs(uid string) map[string]string
	Close()
}
