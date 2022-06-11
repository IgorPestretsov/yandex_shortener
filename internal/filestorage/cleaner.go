package filestorage

type RecordToDelete struct {
	urlID  string
	userID string
}
type Cleaner struct {
	UserDeleteRequests chan map[string]string
	storage            *Storage
	q                  chan bool
}

func NewCleaner(s *Storage) Cleaner {
	reqs := make(chan map[string]string, 100)
	return Cleaner{UserDeleteRequests: reqs, storage: s}
}

func (c *Cleaner) Run(quit chan bool) {
	c.q = quit
	queue := c.fillQueue()
	go func() {
		for {
			select {
			case batch := <-queue:
				c.storage.DeleteRecord(batch)
			case <-c.q:
				return
			}
		}

	}()
}
func (c *Cleaner) fillQueue() chan RecordToDelete {
	out := make(chan RecordToDelete)
	go func() {
		for {
			select {
			case userReq := <-c.UserDeleteRequests:
				for key, uuid := range userReq {
					out <- RecordToDelete{urlID: key, userID: uuid}
				}
			case <-c.q:
				close(out)
				return
			}
		}
	}()
	return out
}
