package storagememory

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestNewShortURLStore(t *testing.T) {
	t.Run("should success create tests", func(t *testing.T) {
		got := NewShortURLStore()

		if got == nil {
			t.Errorf("NewShortURLStore() = %v", got)
		}
	})
}

func Test_shortURLRepository_Add(t *testing.T) {
	t.Run("should success add", func(t *testing.T) {
		storeMap := map[string]*core.ShortURL{}

		s := &shortURLRepository{
			store: storeMap,
			mutex: &sync.RWMutex{},
		}
		shortURL := &core.ShortURL{
			ID:        "1",
			URL:       "https://vk.com",
			UserID:    "1",
			IsDeleted: false,
		}

		if err := s.Add(context.Background(), shortURL); err != nil {
			t.Errorf("shortURLRepository.Add() error = %v", err)
		}

		if !reflect.DeepEqual(shortURL, storeMap["1"]) {
			t.Errorf("shortURLRepository.GetByID() got = %v", shortURL)
		}
	})
}

func Test_shortURLRepository_GetByID(t *testing.T) {
	storeMap := map[string]*core.ShortURL{
		"1": {
			ID:        "1",
			URL:       "https://vk.com",
			UserID:    "1",
			IsDeleted: false,
		},
		"2": {
			ID:        "2",
			URL:       "https://vk.com/2",
			UserID:    "1",
			IsDeleted: false,
		},
	}

	s := &shortURLRepository{
		store: storeMap,
		mutex: &sync.RWMutex{},
	}

	t.Run("should success get by id", func(t *testing.T) {
		got, got1 := s.GetByID(context.Background(), "1")
		if !reflect.DeepEqual(got, storeMap["1"]) {
			t.Errorf("shortURLRepository.GetByID() got = %v", got)
		}

		if got1 != true {
			t.Errorf("shortURLRepository.GetByID() got1 = %v", got1)
		}
	})

	t.Run("should not found", func(t *testing.T) {
		_, got1 := s.GetByID(context.Background(), "3")
		if got1 != false {
			t.Errorf("shortURLRepository.GetByID() got1 = %v", got1)
		}
	})
}

func Test_shortURLRepository_AllByUserID(t *testing.T) {
	storeMap := map[string]*core.ShortURL{
		"1": {
			ID:        "1",
			URL:       "https://vk.com",
			UserID:    "1",
			IsDeleted: false,
		},
		"2": {
			ID:        "2",
			URL:       "https://vk.com/2",
			UserID:    "1",
			IsDeleted: false,
		},
		"3": {
			ID:        "3",
			URL:       "https://vk.com/3",
			UserID:    "3",
			IsDeleted: false,
		},
		"4": {
			ID:        "4",
			URL:       "https://vk.com/4",
			UserID:    "4",
			IsDeleted: false,
		},
	}

	s := &shortURLRepository{
		store: storeMap,
		mutex: &sync.RWMutex{},
	}

	t.Run("should success return by user", func(t *testing.T) {
		got, err := s.AllByUserID(context.Background(), "1")

		if err != nil {
			t.Errorf("shortURLRepository.AllByUserID() error = %v", err)
			return
		}

		assert.Equal(t, 2, len(got))

		if !reflect.DeepEqual(got[0], storeMap["1"]) {
			t.Errorf("shortURLRepository.AllByUserID() = %v with %v", got[0], storeMap["1"])
		}

		if !reflect.DeepEqual(got[1], storeMap["2"]) {
			t.Errorf("shortURLRepository.AllByUserID() = %v with %v", got[0], storeMap["1"])
		}
	})

	t.Run("should return empty slice", func(t *testing.T) {
		got, err := s.AllByUserID(context.Background(), "5")

		if err != nil {
			t.Errorf("shortURLRepository.AllByUserID() error = %v", err)
			return
		}

		assert.Equal(t, 0, len(got))
	})
}

func Test_shortURLRepository_DeleteURLsUserByIds(t *testing.T) {
	t.Run("should success delete URLs by ids", func(t *testing.T) {
		storeMap := map[string]*core.ShortURL{
			"1": {
				ID:        "1",
				URL:       "https://vk.com",
				UserID:    "1",
				IsDeleted: false,
			},
			"2": {
				ID:        "2",
				URL:       "https://vk.com/2",
				UserID:    "1",
				IsDeleted: false,
			},
			"3": {
				ID:        "3",
				URL:       "https://vk.com/3",
				UserID:    "3",
				IsDeleted: false,
			},
			"4": {
				ID:        "4",
				URL:       "https://vk.com/4",
				UserID:    "4",
				IsDeleted: false,
			},
		}

		s := &shortURLRepository{
			store: storeMap,
			mutex: &sync.RWMutex{},
		}

		if err := s.DeleteURLsUserByIds(context.Background(), "1", []string{"1", "2", "3"}); err != nil {
			t.Errorf("shortURLRepository.Add() error = %v", err)
		}

		if storeMap["1"].IsDeleted != true {
			t.Errorf("shortURLRepository.DeleteURLsUserByIds() got = %v", storeMap["1"])
		}
		if storeMap["2"].IsDeleted != true {
			t.Errorf("shortURLRepository.DeleteURLsUserByIds() got = %v", storeMap["2"])
		}
		if storeMap["3"].IsDeleted != false {
			t.Errorf("shortURLRepository.DeleteURLsUserByIds() got = %v", storeMap["3"])
		}
		if storeMap["4"].IsDeleted != false {
			t.Errorf("shortURLRepository.DeleteURLsUserByIds() got = %v", storeMap["4"])
		}
	})
}

//func Test_shortURLRepository_CreateBatchWithContext(t *testing.T) {
//	storeMap := map[string]*core.ShortURL{}
//
//	s := &shortURLRepository{
//		store: storeMap,
//		mutex: &sync.RWMutex{},
//	}
//	t.Run(tt.name, func(t *testing.T) {
//		s := &shortURLRepository{
//			store: tt.fields.store,
//			mutex: tt.fields.mutex,
//		}
//
//		if err := s.CreateBatch(tt.args.in0, tt.args.shortURLs); (err != nil) != tt.wantErr {
//			t.Errorf("shortURLRepository.CreateBatch() error = %v, wantErr %v", err, tt.wantErr)
//		}
//	})
//}
