package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/techjosec/url-shortener-app/api"
	"github.com/techjosec/url-shortener-app/repository/redis"
	"github.com/techjosec/url-shortener-app/shortener"
)

func main() {

	log.Printf("Starting URL-SHORTENER-APP ...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Fatal Error loading .env file")
	}

	redisRepository := getRedisRepository()
	service := shortener.NewRedirectService(redisRepository)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/api/redirect/{code}", handler.Get)
	router.Post("/api/redirect/", handler.Post)

	errors := make(chan error, 2)
	go func() {
		port := httpPort()

		log.Printf("Server listening on port %s \n", port)
		errors <- http.ListenAndServe(port, router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errors <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errors)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("APP_PORT") != "" {
		port = os.Getenv("APP_PORT")
	}

	return fmt.Sprintf(":%s", port)
}

func getRedisRepository() shortener.RedirectRepository {

	if os.Getenv("REDIS_URL") == "" {

		log.Fatal("Fatal Error: Missing REDIS_URL environment variable")

	}

	redisUrl := os.Getenv("REDIS_URL")
	repo, err := redis.NewRedisRepository(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	return repo

}
