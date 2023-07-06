package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robinmuhia/GolangProjects/rssAggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main(){
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not connected")
	}	
	dbURL := os.Getenv("DB_URL")
	if dbURL== ""{
		log.Fatal("DB is not connected")
	}	
	conn,err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Fatal("Can't connect to databse",err)
	}
	db := database.New(conn)
	apiConfig := apiConfig{
		DB: db,
	}
	go startScraping(db,10,time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{		
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	  }))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handleErr)
	v1Router.Post("/users",apiConfig.handlerCreateUser)
	v1Router.Get("/users",apiConfig.middlewareAuth(apiConfig.handlerGetUser))
	v1Router.Post("/feeds",apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds",apiConfig.handlerGetFeeds)
	v1Router.Post("/feed_follows",apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows",apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowId}",apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollow))
	v1Router.Get("/posts",apiConfig.middlewareAuth(apiConfig.handlerGetPostsForUser))
	router.Mount("/v1",v1Router)
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}
	log.Printf("Server starting on port %v",portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Port is listening on PORT: %v",portString)
}