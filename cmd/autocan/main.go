package main

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	indexRequest      int64 = 0
	countWorker             = 350
	endpointAPICreate       = "http://localhost:8080/api/shorten/"
	client                  = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10000,
			MaxIdleConns:        10000,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,
			DisableKeepAlives:   false,
			ForceAttemptHTTP2:   false,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tCtx, tCancel := context.WithTimeout(ctx, 20*time.Second)
	defer tCancel()

	g, gCtx := errgroup.WithContext(tCtx)

	// Для нагружалки сделать более тяжелую и больше го рутин
	for i := 0; i < countWorker; i++ {
		g.Go(func() error {
			for {
				select {
				case <-gCtx.Done():
					//log.Println("close")
					return nil
				default:
					res2 := fmt.Sprintf(`{"url": "https://ya.ru/%v"}`, atomic.AddInt64(&indexRequest, 1))

					req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointAPICreate, bytes.NewBufferString(res2))
					req.Header.Set("Content-Type", "application/json")

					if err != nil {
						return err
					}

					resp, err := client.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()
					_, err = io.ReadAll(resp.Body)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	go func() {
		tiker := time.NewTicker(time.Second)
		defer tiker.Stop()

		for {
			select {
			case <-tiker.C:
				log.Println("count request", atomic.LoadInt64(&indexRequest))
			case <-gCtx.Done():
				tiker.Stop()
				log.Println("Stop timer")
				return
			}
		}
	}()

	err := g.Wait()

	if err != nil {
		log.Println("will error", err)

		return
	}

	log.Println("count request", atomic.LoadInt64(&indexRequest))
	log.Println("Success result")

	return
}
