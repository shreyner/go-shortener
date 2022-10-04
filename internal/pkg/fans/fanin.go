package fans

import "sync"

func fanIn(inputChs []chan *FanDeleteJob) chan *FanDeleteJob {
	outCh := make(chan *FanDeleteJob)

	go func() {
		wg := sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)
			go func(inputCh chan *FanDeleteJob) {
				defer wg.Done()
				for job := range inputCh {
					outCh <- job
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}
