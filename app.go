package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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
		if errors.Is(err, models.ErrEmailExist) {
			app.badErrorResponse(w, "Company with this email already exists")
			return
		}
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

	// var filePath string

	com, err := app.VerifyAPIKey(r.Header.Get("Authorization"))
	if err != nil {
		if errors.Is(err, models.ErrAPiKey) {
			app.badErrorResponse(w, "Invalid API key provided")
			return
		}
		app.serverErrorResponse(w, err)
		return
	}

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 5 MB files
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		app.badErrorResponse(w, toJson{"error": "File too large must not be more than 5MB"})
		return
	}

	file, header, err := r.FormFile("document")

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	if ext != ".pdf" && ext != ".docx" && ext != ".csv" && ext != ".txt" {
		app.serverErrorResponse(w, err)
		return
	}

	reader, err := io.ReadAll(file)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	content, err := GetContentFromFile(reader, ext)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	emb, err := app.GenerateEmbeddings(content)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	doc.ID = env.GetID()
	// doc.DocumentPath = filePath
	doc.CompanyID = com.ID
	doc.Content = content
	doc.Embedding = EmbeddingToString(emb)

	err = app.model.Model.AddDocument(doc)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	app.writeResponse(w, http.StatusCreated, toJson{"message": "Document added successfully"})
}

// HandleChat handles the chat endpoint
func (app *application) HandleChat(w http.ResponseWriter, r *http.Request) {

	prompt := `
	You are an AI-powered FAQ assistant for [Company Name], designed to provide accurate and helpful responses based only on the official company documentation. Your goal is to assist users in a friendly, professional, and personalized manner while staying strictly within the provided information.

	Guidelines for Responses:
	Accuracy & Relevance

	Answer questions directly and concisely using only the provided company information.
	If a question includes multiple parts, address each one separately (if the information is available).
	Do not guess, infer, or fabricate information.
	 Personalization & Engagement

	Use the user's name (if available) or personalized pronouns to create a warm, engaging experience.
	Maintain a friendly, helpful, and professional tone that reflects the company’s brand.
	Make responses clear and structured (use bullet points for lists when appropriate).
	 Handling Missing Information

	If the answer isn't available in the company documentation, respond with:
	"I don't have that specific information. Please contact our support team at support@[company-domain].com for assistance."
	If a question is unclear, politely ask for clarification instead of assuming.
	 Example Responses for Different Scenarios:

	User asks about a service:
	"Great question, Alex! [Company Name] offers [brief description of the service]. If you need more details, let me know!"
	User asks about an unavailable topic:
	"I don't have that specific information, but you can reach out to our support team at support@[company-domain].com for assistance."
	User asks a multi-part question:
	"Good question, Sarah! Here's a breakdown:
	[Feature 1]: [Explanation]
	[Feature 2]: [Explanation]
	Let me know if you'd like further details!"_
	Your priority is to deliver fast, accurate, and engaging responses that enhance the user’s experience while representing [Company Name] professionally
	`
	// process query request coming from telex
	var Query struct {
		Message  string `json:"message"`
		Settings struct {
			Authorization string `json:"Authorization"`
			Type          string `json:"type"`
			Default       string `json:"default"`
		} `json:"settings"`
	}

	err := ReadJson(r, &Query)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	log.Printf("Query: %+v\n", Query)
	app.writeResponse(w, http.StatusOK, toJson{"message": "Chat response generated successfully", "response": Query.Message})

	com, err := app.model.Model.GetAPIKey(Query.Settings.Authorization)
	if err != nil {
		if errors.Is(err, models.ErrAPiKey) {
			app.badErrorResponse(w, "Invalid API key provided")
			return
		}
		app.serverErrorResponse(w, err)
		return
	}

	// com, err := app.VerifyAPIKey(r.Header.Get("Authorization"))

	// if err != nil {
	// 	if errors.Is(err, models.ErrAPiKey) {
	// 		app.badErrorResponse(w, "Invalid API key provided")
	// 		return
	// 	}
	// 	app.serverErrorResponse(w, err)
	// }

	queryEmbedding, err := app.GenerateEmbeddings(Query.Message)
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
			{Role: "user", Content: fmt.Sprintf("Here is what we know about the company: %s\nUser Question: %s", knowledge, Query.Message)},
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
				"created_at": "2025-02-22",
				"updated_at": "2025-02-22",
			},
			"descriptions": map[string]string{
				"app_name":         "chatOrg",
				"app_description":  "The AI FAQ Chatbot helps organizations provide instant, accurate answers to customer inquiries using their own data. By integrating documents and manual inputs, it automates support, reduces response time, and enhances customer satisfaction effortlessly.",
				"app_logo":         "https://as2.ftcdn.net/v2/jpg/12/13/39/93/1000_F_1213399398_I4M3xm84LUZSAuI5llmxa2IPIBi34Mow.jpg",
				"app_url":          "https://unchanged-tawnya-hng-c6a8014b.koyeb.app",
				"background_color": "#fff",
			},
			"is_active":            true,
			"integration_type":     "modifier",
			"integration_category": "Communication & Collaboration",
			"key_features": []string{
				"Smart Responses",
				"Document Integration",
				"Manual Input",
				"AI-Powered Search",
			},
			"author": "Darasimi",
			"settings": []map[string]interface{}{
				{
					"label":    "Authorization",
					"type":     "text",
					"required": true,
					"default":  "fe0e09ecda000917df2eda0d87208a92",
				},
			},
			"target_url": "https://unchanged-tawnya-hng-c6a8014b.koyeb.app/api/v1/chat",
			"tick_url":   "https://unchanged-tawnya-hng-c6a8014b.koyeb.app/api/v1/tick",
		},
	}
	app.writeResponse(w, http.StatusOK, data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("pong")
}
