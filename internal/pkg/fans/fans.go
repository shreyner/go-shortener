package fans

import (
	"context"
	"github.com/shreyner/go-shortener/internal/repositories"
	"go.uber.org/zap"
)

type FanDeleteJob struct {
	UserID string
	URLIDs []string
}

type FansShortService struct {
	inputCh           chan *FanDeleteJob
	shorterRepository repositories.ShortURLRepository
}

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

func (s *FansShortService) Add(userID string, URLIDs []string) {
	job := &FanDeleteJob{
		UserID: userID,
		URLIDs: URLIDs,
	}

	s.inputCh <- job
}

func (s *FansShortService) Close() {
	close(s.inputCh)
}
