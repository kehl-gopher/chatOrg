package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	mux := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      routers(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	err := mux.ListenAndServe()

	log.Fatal(err)
}
