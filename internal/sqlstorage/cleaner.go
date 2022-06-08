package sqlstorage

import (
	"sync"
)

const workersCount = 10

type urlToDelete map[string]string

type RecordToDelete struct {
	urlID  string
	userID string
}

type Cleaner struct {
	UserDeleteRequests chan map[string]string
	storage            *Storage
}

func NewCleaner(s *Storage) Cleaner {
	reqs := make(chan map[string]string, 100)
	return Cleaner{UserDeleteRequests: reqs, storage: s}
}

func (c *Cleaner) Run() {
	queue := c.fillQueue()
	fanOutChs := c.fanOut(queue, workersCount)
	workerChs := make([]chan RecordToDelete, 0, workersCount)

	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan RecordToDelete)
		c.newWorker(fanOutCh, workerCh)
		workerChs = append(workerChs, workerCh)
	}
	toDeleteChn := c.fanIn(workerChs...)
	go func() {
		var batchToDelete []RecordToDelete
		for r := range toDeleteChn {
			batchToDelete = append(batchToDelete, r)
			if len(batchToDelete) > 0 {
				c.storage.DeleteRecords(batchToDelete)
			}
		}

	}()

}

func (c *Cleaner) fillQueue() chan RecordToDelete {
	out := make(chan RecordToDelete)
	go func() {
		for {
			for userReq := range c.UserDeleteRequests {
				for key, uuid := range userReq {
					out <- RecordToDelete{urlID: key, userID: uuid}
				}
			}

		}
		close(out)
	}()
	return out
}
func (c *Cleaner) fanOut(inputCh chan RecordToDelete, n int) []chan RecordToDelete {
	chs := make([]chan RecordToDelete, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan RecordToDelete)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan RecordToDelete) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			batch, ok := <-inputCh
			if !ok {
				return
			}

			ch := chs[i]
			ch <- batch
		}
	}()
	return chs

}
func (c *Cleaner) newWorker(input, out chan RecordToDelete) {
	go func() {
		for record := range input {
			if c.storage.CheckReqToDelete(record) {
				out <- record

			}
		}
		close(out)
	}()
}
func (c *Cleaner) fanIn(inputChs ...chan RecordToDelete) chan RecordToDelete {
	outCh := make(chan RecordToDelete)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan RecordToDelete) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}
