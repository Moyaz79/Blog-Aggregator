package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Moyaz79/Blog-Aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

//struct that can be used on function as method
type apiConfig struct {
	DB *database.Queries
}

func main() {


	// library for loading the env file
	godotenv.Load(".env")

	//getting the port from the env file by us the os.Getenv
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	//connection to database using postgres and the url
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	db := database.New(conn)
	apiCfg := apiConfig {
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

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
	v1Router.Post("/users", apiCfg.createUserHandler)
	//this will convert the getUserHandler to the standard http handler 
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.getUserHandler))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	v1Router.Get("/feeds", apiCfg.getFeedHandler)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.createFeedFollowsHandler))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.getFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.deleteFeedFollowsHandler))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.getPostsForUserHandler))


	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT:", portString)

}
