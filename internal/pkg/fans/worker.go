package fans

func worker(inputCh, outCh chan *FanDeleteJob) {
	go func() {
		for job := range inputCh {
			outCh <- job
		}

		close(outCh)
	}()
}
