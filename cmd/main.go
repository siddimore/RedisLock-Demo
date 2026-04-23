package main

import (
	"distributed_lock_demo/internal"
	"distributed_lock_demo/rideshare"
	"log"
	"net/http"
)

func main() {
	rclient := internal.NewRedisClient()

	http.HandleFunc("/accept", func(w http.ResponseWriter, r *http.Request) {
		rideshare.AcceptRide(w, r, rclient)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
