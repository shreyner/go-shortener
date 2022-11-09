// Package fans implementation fanOut/fanIn pattern for concurrent delete urls
package fans

import (
	"context"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/repositories"
)

// FanDeleteJob job include all field for delete sortn URL
type FanDeleteJob struct {
	UserID string
	URLIDs []string
}

// FansShortService business logic for delte urls
type FansShortService struct {
	inputCh           chan *FanDeleteJob
	shorterRepository repositories.ShortURLRepository
}

// NewFansShortService create worker for bachground works
func NewFansShortService(log *zap.Logger, rep repositories.ShortURLRepository, workerCount int) *FansShortService {
	inputCh := make(chan *FanDeleteJob)

	fansShortService := FansShortService{
		inputCh:           inputCh,
		shorterRepository: rep,
	}

	fanOutChs := fanOut(inputCh, workerCount)
	workerOutChs := make([]chan *FanDeleteJob, 0, workerCount)
	for _, fanOutCh := range fanOutChs {
		workerOutCh := make(chan *FanDeleteJob)
		worker(fanOutCh, workerOutCh)
		workerOutChs = append(workerOutChs, workerOutCh)
	}

	outCh := fanIn(workerOutChs)

	go func(outCh chan *FanDeleteJob) {
		for job := range outCh {
			if err := rep.DeleteURLsUserByIds(context.Background(), job.UserID, job.URLIDs); err != nil {
				log.Error("error when delete urls for user", zap.String("userID", job.UserID), zap.Error(err))
			}
		}
	}(outCh)

	return &fansShortService
}

// Add new job in queue for delete
func (s *FansShortService) Add(userID string, URLIDs []string) {
	job := &FanDeleteJob{
		UserID: userID,
		URLIDs: URLIDs,
	}

	s.inputCh <- job
}

// Close queue and stop process
func (s *FansShortService) Close() {
	close(s.inputCh)
}
