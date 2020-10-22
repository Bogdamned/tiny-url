package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const port = ":8000"
const address = "127.0.0.1" + port

type Server struct {
	Handlers Handler
}

func run(h *Handler) {
	r := mux.NewRouter()

	//h := Handler{NewLocalCache()}
	r.HandleFunc("/{tiny}", h.TinyRedirect).Methods("GET")
	r.HandleFunc("/", h.Tinify).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	gracefulfShutdown(srv)
}

func gracefulfShutdown(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
