package main

import (
	"log"
	"net/http"

	"github.com/shreyner/go-shortener/internal/handler"
)

var listenServerAddress = ":8080"

var mapShortedUrls = map[string]string{
	"yand": "https://yandex.ru",
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.IndexHandler(mapShortedUrls))

	log.Printf("Start start on %s", listenServerAddress)
	log.Fatalln(http.ListenAndServe(listenServerAddress, mux))
}
