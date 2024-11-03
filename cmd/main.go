package main

import (
	"net/http"
	"os"
	"strconv"

	"CustomerCRUD/pkg/repository"
	"CustomerCRUD/server"
	"CustomerCRUD/utils"

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

	dbURL := os.Getenv("DATABASE_URL")
	if err := utils.RunMigrations(dbURL); err != nil {
		log.Fatal("error running db migrations", err)
	}

	dbRepo := repository.NewCustomerRepository(db)

	srv := server.NewServer(dbRepo)
	srv.SetupRoutes()

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Router))
}
