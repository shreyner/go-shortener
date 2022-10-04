package fans

func fanOut(input chan *FanDeleteJob, n int) []chan *FanDeleteJob {
	chs := make([]chan *FanDeleteJob, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan *FanDeleteJob)

		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *FanDeleteJob) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			num, ok := <-input

			if !ok {
				return
			}
			ch := chs[i]
			ch <- num
		}
	}()
	return chs
}
