package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			bids.CreateNewBid()
		}
	}()

	runServer(mux)
}

func runServer(mux *http.ServeMux) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("Error to start the server: %s", err)
		}
	}()
	log.Println("The server started ...")

	<-stop
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Println("Error when shutting down the server:", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
