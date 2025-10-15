package infra

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error .env file")
	}

}
	