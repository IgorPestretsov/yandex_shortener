package storage

type Channels struct {
	KeyChannel        chan string
	FullLinkChannel   chan string
	LinksPairsChannel chan [2]string
}

type Storage struct {
	Channels Channels
	Storage  map[string]string
}

func (s *Storage) Init() {
	s.Channels.LinksPairsChannel = make(chan [2]string)
	s.Channels.KeyChannel = make(chan string)
	s.Channels.FullLinkChannel = make(chan string)
	s.Storage = make(map[string]string)
}

func (s *Storage) Run() {
	s.Init()
	go s.SaveLinksPair()
	go s.LoadLinksPair()
}

func (s *Storage) LoadLinksPair() {
	for {
		key := <-s.Channels.KeyChannel
		FullLink := s.Storage[key]
		s.Channels.FullLinkChannel <- FullLink
	}
}

func (s *Storage) SaveLinksPair() {
	for {
		LinksPair := <-s.Channels.LinksPairsChannel
		s.Storage[LinksPair[1]] = LinksPair[0]
	}
}
