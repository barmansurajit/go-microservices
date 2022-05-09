package main

import (
	"context"
	handlers "github.com/barmansurajit/go-microservices/product-api/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	serveMux := http.NewServeMux()

	// handler associations
	productHandler := handlers.NewProducts(l)
	serveMux.Handle("/", productHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Println(err)
			os.Exit(1)
		}
	}()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	s := <-signalChannel
	log.Println("Received terminate, graceful shutdown", s)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
	os.Exit(0)
}
