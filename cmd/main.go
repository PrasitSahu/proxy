package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	conf "github.com/PrasitSahu/proxy/internal"
)

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		res.Write([]byte("Hello World!"))
	})

	server := &http.Server{
		Addr: ":" + conf.Config.Port,
		Handler: mux,
		MaxHeaderBytes: 1 << 20,
	}

	startServer(server)
	listenToSignal()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	shutdownServer(ctx, server)
}

func startServer(server *http.Server){
	log.Printf("starting server on port %s", conf.Config.Port)
	go func(){
		if err := server.ListenAndServe(); err != nil {
			log.Println("server shutting down: ", err)
		}
	}()
}

func shutdownServer(ctx context.Context, server *http.Server){
	if err := server.Shutdown(ctx); err != nil {
		log.Println("shutting down server...")
	}
}

func listenToSignal(){
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	<- sig
}