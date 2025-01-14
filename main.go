package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {

	WEBHOOK_URL := flag.String("webhook_url", "http://localhost:3000/webhook", "where do you want your emitted event to go?")
	flag.Parse()

	log.Printf("Sending events on webhook URL %s\n", *WEBHOOK_URL)

	mux := http.NewServeMux()
	timer := CreateTimer(*WEBHOOK_URL)

	// System goroutine;
	timer.SysWg.Add(1)

	// Handle when a user tries to register a timeout
	// Spawn the goroutine
	mux.HandleFunc("POST /register", timer.RegisterHandler)
	mux.HandleFunc("POST /cancel", timer.CancelHandler)
	mux.HandleFunc("POST /remaining", timer.RemainingHandler)
	mux.HandleFunc("POST /extend", timer.ExtendHandler)

	// placeholder for the webhook
	mux.HandleFunc("POST /webhook", timer.WebhookTest)

	var systemwg sync.WaitGroup
	systemwg.Add(2)

	// Spawn server goroutine
	go func() {
		log.Println("Starting server ...")
		err := http.ListenAndServe(":3000", mux)
		if err != nil {
			log.Println(fmt.Errorf("Error in server -> %v", err))
			defer systemwg.Done()
		}
	}()

	// Spawn timer goroutine
	go func() {
		defer systemwg.Done()
		timer.SysWg.Wait()
	}()

	systemwg.Wait()
}
