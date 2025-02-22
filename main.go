package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"telex-chat/internal/env"
	"telex-chat/internal/models"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
)

type application struct {
	openai *openai.Client
	model  *models.AppModel
	cld    *cloudinary.Cloudinary
	env    string
}

func main() {
	env := env.DotEnv("ENV")
	db, err := initDB(env)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// get environment variable
	cld, err := instCloudinary()
	if err != nil {
		log.Fatal(err)
	}
	app := &application{model: models.NewAppModel(db), openai: instOpenAIClient(), cld: cld, env: env}

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

// initDB initializes the database connection
func initDB(envi string) (*sql.DB, error) {
	var db_str string

	if envi == "DEV" {
		log.Println("Running in development mode")
		db_str = env.DotEnv("DB_URL")
	} else {
		log.Println("Running in production mode")
		db_str = env.DotEnv("PROD_DB_URL")
	}
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

// instOpenAIClient initializes the openai client
func instOpenAIClient() *openai.Client {
	apiKey :=
		env.DotEnv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return client
}

func instCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := env.DotEnv("CLOUDINARY_CLOUD_NAME")
	apiKey := env.DotEnv("CLOUDINARY_API_KEY")
	apiSecret := env.DotEnv("CLOUDINARY_API_SECRET")
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cld, nil
}
