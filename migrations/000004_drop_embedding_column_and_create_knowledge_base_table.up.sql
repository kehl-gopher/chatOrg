DROP TABLE faq;


-- ALTER TABLE documents DROP COLUMN embedding;
ALTER TABLE about DROP COLUMN embedding;

CREATE TABLE knowledge_base (
    id TEXT UNIQUE PRIMARY KEY,
    content TEXT,
    embedding vector(1536),
    company_id TEXT,
    FOREIGN KEY (company_id) REFERENCES company (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);