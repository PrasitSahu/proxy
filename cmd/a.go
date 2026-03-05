package main

import (
	"log"
	"os"
	"strings"

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
	if len(strings.TrimSpace(conf.Config.Port)) == 0 {
		log.Fatal("variable 'PORT' is unset")
	}
}