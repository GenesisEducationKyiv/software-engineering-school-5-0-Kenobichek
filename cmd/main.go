package main

import (
	"log"
	"net/http"
	"os"

	"Weather-Forecast-API/internal"
	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/scheduler"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file.")
	}
}

func main() {
	db.Init()
	db.RunMigrations(db.DataBase)

	go scheduler.StartScheduler()

	router := chi.NewRouter()
	internal.RegisterRoutes(router)

	port := os.Getenv("PORT")

	log.Println("Server started at :" + port)
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		return
	}
}
