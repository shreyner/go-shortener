package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/shreyner/go-shortener/proto"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	authClint := pb.NewAuthClient(conn)
	shortenerClient := pb.NewShortenerClient(conn)

	ctx := context.Background()

	log.Println("Create token")
	getTokenResponse, err := authClint.GetToken(ctx, &pb.Empty{})

	if err != nil {
		log.Println("Error get token", err)
		return
	}

	log.Println("Token: ", getTokenResponse.Token)

	md := metadata.New(map[string]string{
		"token": getTokenResponse.Token,
	})

	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	urls := []*pb.CreateShortRequest{
		{Url: fmt.Sprintf("https://ya.ru/%v", random.Int())},
		{Url: fmt.Sprintf("https://ya.ru/%v", random.Int())},
	}

	log.Println("Create single short")
	for _, url := range urls {
		resp, createErr := shortenerClient.CreateShort(ctxWithToken, url)
		if createErr != nil {
			log.Println("Error when create", createErr)
			return
		}

		if resp.Error != "" {
			log.Println("result error", resp.Error)
			return
		}

		log.Printf("Was create url: %v, id: %v\n", url.Url, resp.Id)
	}

	urlsBatch := []*pb.CreateBatchShortRequest_URLs{
		{Url: fmt.Sprintf("https://vk.com/%v", random.Int()), CorrelationId: strconv.FormatInt(int64(random.Int()), 10)},
		{Url: fmt.Sprintf("https://vk.com/%v", random.Int()), CorrelationId: strconv.FormatInt(int64(random.Int()), 10)},
	}

	batchShortRequest := pb.CreateBatchShortRequest{
		Urls: urlsBatch,
	}

	log.Println("Create batch short")
	respCreateBatch, err := shortenerClient.CreateBatchShort(ctxWithToken, &batchShortRequest)

	if err != nil {
		log.Println("Error when create", err)
		return
	}

	if respCreateBatch.Error != "" {
		log.Println("result error", respCreateBatch.Error)
		return
	}

	for _, url := range respCreateBatch.Urls {
		log.Printf("Create short url id: %v, CorrelationId: %v\n", url.Id, url.CorrelationId)
	}

	log.Println("Get list short")
	listUserURLS, err := shortenerClient.ListUserURLs(ctxWithToken, &pb.ListUserURLsRequest{})

	if err != nil {
		log.Println("Error when create", err)
		return
	}

	for _, url := range listUserURLS.Urls {
		log.Printf("user url id: %v, originURL: %v\n", url.Id, url.OriginalURL)
	}

	ids := make([]string, len(listUserURLS.Urls))

	for i, url := range listUserURLS.Urls {
		ids[i] = url.Id
	}

	log.Printf("Will be delete by ids: %v", ids)
	_, err = shortenerClient.DeleteByIDs(ctxWithToken, &pb.DeleteByIDsRequest{Ids: ids})

	if err != nil {
		log.Println("Error when delete all urls", err)
		return
	}

	time.Sleep(10 * time.Second)

	log.Println("Get list short after all delete")
	list, err := shortenerClient.ListUserURLs(ctxWithToken, &pb.ListUserURLsRequest{})

	if err != nil {
		log.Println("Error when create", err)
		return
	}

	log.Printf("len urls for user: %v", len(list.Urls))
}
