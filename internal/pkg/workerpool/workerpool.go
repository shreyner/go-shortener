package workerpool

import (
	"runtime"
	"sync"

	"go.uber.org/zap"
)

type JobDeleteURLs struct {
	UserID string
	URLIDs []string
}

type Queue struct {
	arr  []*JobDeleteURLs
	mu   sync.Mutex
	cond *sync.Cond
}

func NewQueue() *Queue {
	q := Queue{}
	q.cond = sync.NewCond(&q.mu)

	return &q
}

func (q *Queue) Push(task *JobDeleteURLs) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.arr = append(q.arr, task)
	q.cond.Signal()
}

func (q *Queue) PopWait() *JobDeleteURLs {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.arr) == 0 {
		q.cond.Wait()
	}

	t := q.arr[0]

	q.arr = q.arr[1:]

	return t
}

type JobDeleter func(*JobDeleteURLs) error

type Worker struct {
	id         int
	log        *zap.Logger
	queue      *Queue
	jobDeleter JobDeleter
}

func (w *Worker) Loop(stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
		default:
			t := w.queue.PopWait()

			if err := w.jobDeleter(t); err != nil {
				w.log.Error("error worker", zap.Int("workerID", w.id), zap.Error(err))
			}
		}

	}
}

type WorkerPool struct {
	workers []*Worker
	queue   *Queue

	stopCh chan struct{}
}

func NewWorkerPool(log *zap.Logger, jobDeleter JobDeleter) *WorkerPool {
	wp := WorkerPool{
		queue:   NewQueue(),
		workers: make([]*Worker, 0, runtime.NumCPU()),

		stopCh: make(chan struct{}),
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		wp.workers = append(wp.workers, &Worker{
			id:         i,
			log:        log,
			queue:      wp.queue,
			jobDeleter: jobDeleter,
		})
	}

	for _, worker := range wp.workers {
		go worker.Loop(wp.stopCh)
	}

	return &wp
}

func (wp *WorkerPool) Push(job *JobDeleteURLs) {
	wp.queue.Push(job)
}

func (wp *WorkerPool) Stop() {
	close(wp.stopCh)
}
