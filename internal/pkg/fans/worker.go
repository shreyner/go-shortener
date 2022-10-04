package fans

import "log"

func worker(inputCh, outCh chan *FanDeleteJob) {
	go func() {
		for num := range inputCh {
			log.Println(num)
			outCh <- num
		}

		close(outCh)
	}()
}
