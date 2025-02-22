# AI FAQ Chatbot Documentation

## Overview
The AI FAQ chatbot is designed to assist users by providing relevant answers based on a company's documentation and manually entered information. It utilizes OpenAI embeddings and GPT-based responses to ensure accurate and relevant responses.

## Purpose
The chatbot is part of the Telex integration application, allowing companies to set up an AI-driven FAQ system. By integrating with Telex, businesses can automate responses to common inquiries using their uploaded documents and manually entered data.

## Features
1. **Company Data Sources**
   - Companies can **upload a document** (accepted formats: PDF, TXT, DOCX). The document size must not exceed **10MB**.
   - Companies can **manually enter textual descriptions** of their services.
   - Companies can use **both** data sources for better response accuracy.

**NOTE** PDF uploaded is not yet fully supported

2. **Processing Flow**
   - When a company uploads a document, the system **extracts text** from it.
   - Extracted text and manually entered descriptions are **stored as embeddings** in the database.
   - When a user asks a question, the chatbot uses **OpenAI embeddings** to find the most relevant context and passes it to **GPT** to generate a response.

3. **Technology Stack**
   - **Backend:** Go (with `httprouter` for routing)
   - **Database:** PostgreSQL (with `vector` extension for embeddings)
   - **AI Integration:** OpenAI API (Embeddings + ChatGPT)
   - **Document Processing:** PDF/Text extraction

## API Endpoints

### 1. Create a Company Account
- **Endpoint:** `/api/v1/company`
- **Method:** `POST`
- **Handler:** `app.RegisterCompany`

### 2. Add Company Information (About Section)
- **Endpoint:** `/api/v1/about`
- **Method:** `POST`
- **Headers:**
  ```
  Content-Type: application/json
  Authorization: <API_KEY>
  ```
- **Request Body:**
  ```json
  {
    "about": "Telex gives developers a 90% royalty on their app integration while Telex takes 10%. Telex is affordable, costing only $100 for integration."
  }
  ```
- **Example Curl Command:**
  ```sh
  curl -iX POST "http://localhost:4000/api/v1/about" \
       -H "Content-Type: application/json" \
       -H "Authorization: 32739ab80316e48bca1f10a338a1fca7" \
       -d '{ "about": "Telex gives developers a 90% royalty on their app integration..." }'
  ```
- **Responses:**
  - **Success:** `{ "message": "company info added successfully" }`
  - **Error (Invalid API Key):** `{ "message": "invalid api key" }`

### 3. Upload a Document
- **Endpoint:** `/api/v1/document`
- **Method:** `POST`
- **Handler:** `app.UploadDocument`
- **Request:**
  - Required field: `document`
  - Accepted formats: PDF, TXT, DOCX
  - Max size: 10MB
- **Example Curl Command:**
  ```sh
  curl -iX POST "http://localhost:4000/api/v1/document" \
       -H "Content-Type: application/json" \
       -H "Authorization: 32739ab80316e48bca1f10a338a1fca7" \
       -F "document=@path/to/file.pdf"
  ```
- **Responses:**
  - **Success:** `{ "message": "document uploaded successfully" }`
  - **Error (Invalid API Key):** `{ "message": "invalid api key" }`

### 4. Query the Chatbot
- **Endpoint:** `/api/v1/chat`
- **Method:** `POST`
- **Headers:**
  ```
  Content-Type: application/json
  Authorization: <API_KEY>
  ```
- **Request Body:**
  ```json
  {
    "query": "Tell me about Telex"
  }
  ```
- **Example Curl Command:**
  ```sh
  curl -iX POST "http://localhost:4000/api/v1/chat" \
       -H "Content-Type: application/json" \
       -H "Authorization: 32739ab80316e48bca1f10a338a1fca7" \
       -d '{ "query": "Tell me about Telex" }'
  ```
- **Responses:**
  - **Success:**
    ```json
    {
      "message": "Chat response generated successfully",
      "response": {
        "response": "Telex is a software integration company that simplifies API connectivity for businesses, offering flexible and scalable solutions."
      }
    }
    ```
  - **Error (Invalid API Key):** `{ "message": "invalid api key" }`

## Integration with Telex
To integrate the chatbot with Telex:
1. **Create a Telex account**: [https://telex.im/](https://telex.im/)
2. **Enable ChatOrg**: Navigate to **Settings** > **Enable ChatOrg**.
3. **Manage Apps**: Go to **Manage Apps** and set up your API keys.

## Conclusion
This AI FAQ chatbot enables businesses to automate their customer support by leveraging their own data. By combining document uploads and manual text entries, the chatbot ensures accurate responses using OpenAI's powerful AI models.