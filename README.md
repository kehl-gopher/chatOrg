# AI FAQ Chatbot Documentation

## Overview
chatOrg is a powerful customer support tool that helps businesses automate responses to common inquiries. It is integrated with **Telex**, an AI-powered chatbot system that allows organizations to provide seamless customer experiences using their own documentation and manually entered information.

## Purpose
The chatbot is a  feature of the **Telex** platform, enabling businesses to leverage AI for answering frequently asked questions. Companies can enhance customer interactions by integrating their documentation and textual descriptions into the chatbot system, ensuring precise and relevant responses.

## Features
1. **Company Data Sources**
   - Businesses can **upload a document** (Accepted formats: PDF, TXT, DOCX). The document size must not exceed **10MB**.
   - Companies can **manually enter textual descriptions** of their services for AI reference.
   - Both sources can be combined for enhanced response accuracy.

**⚠️ NOTE:**
- **PDF upload is not fully supported yet.**
- **File upload is not supported currently; documents are processed and stored in the database for embeddings.**

2. **Processing Flow**
   - When a document is uploaded, the system **extracts text** from it.
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

```sh
url -iX POST "http://localhost:4000/api/v1/company" \
     -H "Content-Type: application/json" \
     -d '{"name": "ExampleCorp", "email": "company1234@example.com"}'
```
```json
{"company":{"id":"cusssc1d12354i9649i0",
"name":"ExampleCorp",
"email":"company133234@example.com","api_key":"e4586384fd37ccad237f1f588d475eab84f41013bc731969d20e9cdb12db2a24145e1e46de568eeb"},
"message":"Company registered successfully"}   
```
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
    "about": "Telex allows developers to earn 90% royalties on app integrations while Telex takes 10%. The service is affordable, costing only $100 for integration."
  }
  ```
- **Example Curl Command:**
  ```sh
  curl -X POST "https://api.telex.im/v1/about" \
       -H "Content-Type: application/json" \
       -H "Authorization: 32739ab80316e48bca1f10a338a1fca7" \
       -d '{ "about": "Telex allows developers to earn 90% royalties on app integrations..." }'
  ```
- **Responses:**
  - **Success:** `{ "message": "Company info added successfully" }`
  - **Error (Invalid API Key):** `{ "message": "Invalid API key" }`

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
  curl -X POST "https://api.telex.im/v1/document" \
       -H "Content-Type: application/json" \
       -H "Authorization: 32739ab80316e48bca1f10a338a1fca7" \
       -F "document=@path/to/file.pdf"
  ```
- **Responses:**
  - **Success:** `{ "message": "Document uploaded successfully" }`
  - **Error (Invalid API Key):** `{ "message": "Invalid API key" }`

### 4. Query the Chatbot
- **Endpoint:** `/api/v1/chatorg/query`
- **Method:** `POST`
- **Headers:**
  ```
  Content-Type: application/json
  Authorization: <API_KEY>
  ```
- **Request Body:**
  ```json
  {
    "message": "hi",
    "settings": [
      {
        "Label": "Authorization",
        "Type": "text",
        "Default": "fe0e09ecda000917df2eda0d87208a92",
        "Required": true
      }
    ]
  }
  ```
- **Example Curl Command:**
  ```sh
  curl -X POST "https://api.telex.im/v1/chatorg/query" \
       -H "Content-Type: application/json" \
       -d '{
             "message": "hi",
             "settings": [
               {
                 "Label": "Authorization",
                 "Type": "text",
                 "Default": "fe0e09ecda000917df2eda0d87208a92",
                 "Required": true
               }
             ]
           }'
  ```
- **Responses:**
  - **Success:**
    ```json
    {
      "message": "Chat response generated successfully",
      "response": {
        "response": "Telex is an AI-driven software integration platform that simplifies customer interactions using AI-powered chatbots."
      }
    }
    ```
  - **Error (Invalid API Key):** `{ "message": "Invalid API key" }`

## Integration with Telex ChatOrg
To integrate the chatbot with Telex ChatOrg:
1. **Create a Telex account**: [https://telex.im/](https://telex.im/)
2. **Enable ChatOrg**: Navigate to **Settings** > **Enable ChatOrg**.
3. **Manage Apps**: Go to **Manage Apps** and set up your API keys.

## Conclusion
The **Telex ChatOrg AI FAQ chatbot** allows businesses to automate customer support using AI-powered responses. By leveraging **document uploads and manual text entries**, the chatbot ensures accurate responses based on company-specific data, improving efficiency and customer satisfaction.

