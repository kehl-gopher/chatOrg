# Chat Me

1. **Company Data Sources:**  
   - A company can **upload a document** (PDF, TXT, etc.).  
   - A company can **manually enter text** describing their services.  
   - A company can do **both**.  

2. **Processing Flow:**  
   - When a company uploads a document, we **extract text** from it.  
   - We **store embeddings** for the extracted text and manual input in the database.  
   - When a user asks a question, we use **OpenAI embeddings** to find the most relevant context and pass it to **GPT** to generate a response.  

3. **Technology Stack:**  
   - **Go** with `httprouter` for routing.  
   - **Raw SQL** (PostgreSQL with `vector` extension for embeddings).  
   - **OpenAI API** (Embeddings + ChatGPT).  
   - **Document processing** (PDF/Text extraction).  
