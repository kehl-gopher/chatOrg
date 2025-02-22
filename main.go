package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"telex-chat/internal/env"
	"telex-chat/internal/models"
	"time"

	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
)

type application struct {
	openai *openai.Client
	model  *models.AppModel
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	app := &application{model: models.NewAppModel(db), openai: instOpenAIClient()}

	mux := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      app.routers(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	err = mux.ListenAndServe()

	log.Fatal(err)
}

func initDB() (*sql.DB, error) {
	db_str := env.DotEnv("DB_URL")
	db, err := sql.Open("postgres", db_str)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func instOpenAIClient() *openai.Client {
	apiKey :=
		env.DotEnv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return client
}
