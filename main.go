package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	portStrign := os.Getenv("PORT")
	if portStrign == "" {
		log.Fatal("PORT is not found in the environment")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/err", errHandler)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStrign,
	}

	log.Printf("Server startign on port %v", portStrign)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT:", portStrign)

}
