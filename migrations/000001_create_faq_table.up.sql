CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE company (id TEXT UNIQUE PRIMARY KEY, name TEXT);

CREATE TABLE faq (
    id TEXT UNIQUE PRIMARY KEY,
    questions TEXT NOT NULL,
    answers TEXT NOT NULL,
    company_id TEXT,
    FOREIGN KEY (company_id) REFERENCES company (id) ON DELETE CASCADE
);

CREATE TABLE documents (
    id TEXT UNIQUE PRIMARY KEY,
    doc_path TEXT,
    company_id TEXT,
    FOREIGN KEY (company_id) REFERENCES company (id) ON DELETE CASCADE
);

CREATE TABLE about (
    id TEXT UNIQUE PRIMARY KEY,
    info TEXT,
    company_id TEXT,
    FOREIGN KEY (company_id) REFERENCES company (id) ON DELETE CASCADE
);


ALTER TABLE faq ADD COLUMN embedding vector(1536);
ALTER TABLE about ADD COLUMN embedding vector(1536);




