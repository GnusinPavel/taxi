package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GnusinPavel/taxi/bids"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/request", func(writer http.ResponseWriter, request *http.Request) {
		bid := bids.GetRandom()
		_, err := writer.Write([]byte(bid.Name))
		if err != nil {
			log.Printf("Can't write a response: %s", err)
		}
	})
	mux.HandleFunc("/admin/request", func(writer http.ResponseWriter, request *http.Request) {
		statistics := bids.GetStatistics()
		counter := 0
		for i := range statistics {
			bid := statistics[i]
			if bid.Count > 0 {
				counter++
				_, err := writer.Write([]byte(fmt.Sprintf("%d - %s: %d\n", counter, bid.Name, bid.Count)))
				if err != nil {
					log.Printf("Can't write a response: %s", err)
				}
			}
		}
	})

	ticker := time.NewTicker(time.Millisecond * 200)
	go func() {
		for range ticker.C {
			bids.CreateNewBid()
		}
	}()

	done := make(chan struct{})
	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("не смог запустить сервер: %s", err)
		}
		done <- struct{}{}
	}()
	log.Println("Сервер запущен ...")
	<-done
	ticker.Stop()
}
