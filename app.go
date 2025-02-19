package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// handle FAQ creation
func createFaq(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	json.NewEncoder(w).Encode(body)
}

func webHook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))

	json.NewEncoder(w).Encode(body)
}

func appIntegration(w http.ResponseWriter, r *http.Request) {
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
			"target_url": "https://rfs7htn4-4000.uks1.devtunnels.ms/api/v1/webhook",
			"tick_url":   "https://rfs7htn4-4000.uks1.devtunnels.ms/api/v1/webhook",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
