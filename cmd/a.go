package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	conf "github.com/PrasitSahu/proxy/internal"
	"github.com/joho/godotenv"
)

func init(){
	log.Println("loading env...")
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("failed to load env: ", err)
	}

	conf.Config.Port = os.Getenv("PORT")
	if isEmpty(conf.Config.Port){
		log.Fatal("variable 'PORT' is unset")
	}

	conf.Config.Signature = os.Getenv("SIGNATURE")
	if isEmpty(conf.Config.Signature){
		log.Fatal("variable 'SIGNATURE' is unset")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conf.Config.HttpClient = &http.Client{Timeout: time.Second * 2, Transport: tr}
}

func isEmpty(param string) bool{
	return len(strings.TrimSpace(param)) == 0
}