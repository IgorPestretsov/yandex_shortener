package storage

type Storage interface {
	LoadLinksPair(key string) string
	SaveLinksPair(uid string, key string, link string)
	GetAllUserURLs(uid string) map[string]string
	Close()
}
