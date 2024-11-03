package main

import (
	"os"
	"strconv"

	"CustomerCRUD/pkg/repository"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	isLocalDB := os.Getenv("LOCAL_DB")
	localDBBool, err := strconv.ParseBool(isLocalDB)
	if err != nil {
		log.Fatal("Error parsing local db boolean")
	}

	db, err := repository.GetDB(localDBBool, "DATABASE_URL")
	defer db.Close()
}
