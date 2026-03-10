package main

import (
	"context"
	"crypto/tls"
	// "encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	
	conf "github.com/PrasitSahu/proxy/internal"
	"github.com/PrasitSahu/proxy/internal/api"
)

func main(){
	mux := http.NewServeMux()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Timeout: time.Second * 2, Transport: tr}

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		urlStr := req.URL.Query().Get("url")
		if len(strings.TrimSpace(urlStr)) == 0 {
			http.Error(res, api.ErrNoURL.Error(), http.StatusBadRequest)
			return
		}

		url, err := url.Parse(urlStr)
		if err != nil || url.Host == "" || url.Scheme == "" {
			http.Error(res, api.ErrInvalidURL.Error(), http.StatusBadRequest)
			return
		}

		newRequest, err := http.NewRequestWithContext(
			req.Context(),
			req.Method,
			url.String(),
			req.Body,
		)
		if err != nil {
			http.Error(res, api.ErrReqFail.Error(), http.StatusInternalServerError)
			return
		}

		newRequest.Header = req.Header.Clone()

		resp, err := client.Do(newRequest)
		if err != nil {
			http.Error(res, api.ErrReqFail.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		for k, v := range resp.Header {
			for _, vv := range v {
				res.Header().Add(k, vv)
			}
		}

		res.WriteHeader(resp.StatusCode)
		io.Copy(res, resp.Body)

		return
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