package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"telex-chat/internal/data"
	"telex-chat/internal/env"
	"telex-chat/internal/models"

	"github.com/sashabaranov/go-openai"
)

// RegisterCompany registers a new company
func (app *application) RegisterCompany(w http.ResponseWriter, r *http.Request) {
	company := data.Company{}
	err := ReadJson(r, &company)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	company.ID = env.GetID()
	company.ApiKey, err = data.GenerateSecureAPIKey()
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	comp, err := app.model.Model.AddCompany(company)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	app.writeResponse(w, http.StatusCreated, toJson{"message": "Company registered successfully", "company": comp})
}

// AboutEndpoint handles the about endpoint of the company
func (app *application) AboutEndpoint(w http.ResponseWriter, r *http.Request) {

	var abt data.About
	com, err := app.VerifyAPIKey(r.Header.Get("Authorization"))

	if err != nil {
		if errors.Is(err, models.ErrAPiKey) {
			app.badErrorResponse(w, "Invalid API key provided")
			return
		}
		app.serverErrorResponse(w, err)
		return
	}

	err = ReadJson(r, &abt)
	if err != nil { // check if the json data is valid
		app.badErrorResponse(w, err.Error())
		return
	}

	abt.CompanyID = com.ID
	abt.ID = env.GetID()
	embed, err := app.GenerateEmbeddings(abt.About)
	if err != nil {
		app.serverErrorResponse(w, err.Error)
		return
	}
	abt.Embedding = EmbeddingToString(embed)
	_, err = app.model.Model.AddAbout(abt)
	if err != nil {
		app.serverErrorResponse(w, err.Error())
		return

	}
	app.writeResponse(w, http.StatusCreated, toJson{"message": "company info added successfully"})
}

func (app *application) UploadDocument(w http.ResponseWriter, r *http.Request) {

	var doc data.Document

	com, err := app.VerifyAPIKey(r.Header.Get("Authorization"))
	if err != nil {
		if errors.Is(err, models.ErrAPiKey) {
			app.badErrorResponse(w, "Invalid API key provided")
			return
		}
		app.serverErrorResponse(w, err)
	}

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("document")

	if err != nil {
		http.Error(w, "Invalid file provided", http.StatusBadRequest)
		return
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	if ext != ".pdf" && ext != ".docx" && ext != ".csv" && ext != ".txt" {
		http.Error(w, "Invalid file format provided", http.StatusBadRequest)
		return
	}

	filepath := "uploads/" + header.Filename

	outFile, err := os.Create(filepath)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	doc.ID = env.GetID()
	doc.DocumentPath = filepath
	doc.CompanyID = com.ID

}

// HandleChat handles the chat endpoint
func (app *application) HandleChat(w http.ResponseWriter, r *http.Request) {

	prompt := `
	You are an AI FAQ assistant for a company. Your purpose is to provide helpful, accurate responses about the company based ONLY on the information provided in the company documentation below.
	Guidelines:

	Answer questions directly and concisely using ONLY the provided company information.
	Make sure you personalize query to the user instead of using the word (user's) you can use the user's name. or use a personalized pronouns to make the user feel comfy
	If the answer isn't contained in the company information, respond with: "I don't have that specific information. Please contact the app administrator at support@[company-domain].com for assistance."
	Do not make up or infer information that isn't explicitly stated in the company documentation.
	Maintain a professional, helpful tone that represents the company well.
	If a question is unclear, politely ask for clarification.
	If a question contains multiple parts, address each part if the information is available.
	Format responses for readability when appropriate (bullet points for lists, etc.)
	`
	com, err := app.VerifyAPIKey(r.Header.Get("Authorization"))

	if err != nil {
		if errors.Is(err, models.ErrAPiKey) {
			app.badErrorResponse(w, "Invalid API key provided")
			return
		}
		app.serverErrorResponse(w, err)
	}

	var Query struct {
		Query string `json:"query"`
	}

	err = ReadJson(r, &Query)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	queryEmbedding, err := app.GenerateEmbeddings(Query.Query)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	knowledge, err := app.model.Model.GetMostRelevantKnowledge(com.ID, EmbeddingToString(queryEmbedding))
	if err != nil {
		log.Println(err)
		app.serverErrorResponse(w, err)
		return
	}

	resp, err := app.openai.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: fmt.Sprintf("Here is what we know about the company: %s\nUser Question: %s", knowledge, Query.Query)},
		},
	})

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	Response := struct {
		Response string `json:"response"`
	}{
		Response: resp.Choices[0].Message.Content,
	}
	app.writeResponse(w, http.StatusOK, toJson{"message": "Chat response generated successfully", "response": Response})

}

func (app *application) appIntegration(w http.ResponseWriter, r *http.Request) {
	// integration app settings...
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"date": map[string]string{
				"created_at": "2025-02-19",
				"updated_at": "2025-02-19",
			},
			"descriptions": map[string]string{
				"app_name":         "chatMe",
				"app_description":  "The chatMe is an AI-powered assistant designed to provide instant, accurate responses to frequently asked questions (FAQs) about a company or service. It integrates seamlessly with websites, allowing businesses to automate customer support, \nreduce response times, and enhance user experience.\n\nThis chatbot leverages OpenAI’s GPT models to generate natural and context-aware responses based on preloaded FAQs or dynamically scraped information from a company’s website. It supports real-time learning by fetching and analyzing website content, ensuring up-to-date responses.",
				"app_logo":         "https://i.postimg.cc/Dwr2m6vY/meetme.png",
				"app_url":          "https://rfs7htn4-4000.uks1.devtunnels.ms",
				"background_color": "#fff",
			},
			"is_active":        true,
			"integration_type": "modifier",
			"key_features": []string{
				"Plug-and-Play Integration",
				"Automated Knowledge Extraction",
				"Natural Language Processing (NLP)",
				"Multi-Channel Support",
			},
			"author":               "Darasimi",
			"integration_category": "Communication & Collaboration",
			"settings": []map[string]interface{}{
				{
					"label":    "web_link",
					"type":     "text",
					"required": true,
					"default":  "",
				},
				{
					"label":    "input",
					"type":     "text",
					"required": true,
					"default":  "",
				},
			},
			"target_url": env.DotEnv("TARGET_URL"),
			"tick_url":   env.DotEnv("TICK_URL"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("pong")
}
