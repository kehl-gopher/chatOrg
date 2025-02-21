package models

import (
	"context"
	"database/sql"
	"errors"
	"telex-chat/internal/data"
	"telex-chat/internal/env"
	"time"
)

var ErrAPiKey = errors.New("API key not found")

type NewDB struct {
	Db *sql.DB
}

func NewDBConn(db *sql.DB) *NewDB {
	return &NewDB{Db: db}
}

func (db *NewDB) AddCompany(company data.Company) (*data.Company, error) {

	var com data.Company
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := db.Db.QueryRowContext(ctx,
		"INSERT INTO company (id, name, email, api_key) VALUES ($1, $2, $3, $4) RETURNING id, name, email, api_key",
		company.ID, company.Name,
		company.Email, company.ApiKey,
	).Scan(&com.ID, &com.Name, &com.Email, &com.ApiKey)

	if err != nil {
		return nil, err
	}
	return &com, nil
}

// GetAPIKey retrieves the API key for a company
func (db *NewDB) GetAPIKey(apiKey string) (*data.Company, error) {
	var com data.Company
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := db.Db.QueryRowContext(ctx, "SELECT id, api_key FROM company WHERE api_key = $1",
		apiKey).Scan(&com.ID, &com.ApiKey)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAPiKey
		default:
			return nil, err

		}
	}
	return &com, nil
}

func (db *NewDB) GetCompany(id string) (*data.Company, error) {
	var com data.Company
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := db.Db.QueryRowContext(ctx, "SELECT id, name FROM company WHERE id = $1", id).Scan(&com.ID, &com.Name)

	if err != nil {
		return nil, err
	}
	return &com, nil
}

// func (db *NewDB) AddFAQ(faq data.FAQ) (*data.FAQ, error) {
// 	var faqData data.FAQ
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// embedding to struct
// 	embeddingJson, err := json.Marshal(faq.Embedding)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = db.Db.QueryRowContext(ctx,
// 		"INSERT INTO faq (id, question, answer, embedding,company_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
// 		faq.ID, faq.Question, faq.Answer, embeddingJson, faq.CompanyID).Scan(&faqData.ID)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &faqData, nil
// }

// func (db *NewDB) GetFAQ(id string) (*data.FAQ, error) {
// 	var faq data.FAQ
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	err := db.Db.QueryRowContext(ctx, "SELECT id, question, answer, company_id FROM faq WHERE id = $1", id).Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.CompanyID)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &faq, nil
// }

func (db *NewDB) AddAbout(about data.About) (*data.About, error) {
	var abt data.About
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRowContext(ctx, "INSERT INTO about (id, info, company_id) VALUES ($1, $2, $3) RETURNING id",
		about.ID, about.About, about.CompanyID).Scan(&abt.ID)

	if err != nil {
		return nil, err
	}

	knowledgeBaseId := env.GetID()

	// insert into knowledge base

	_, err = tx.ExecContext(ctx, "INSERT INTO knowledge_base(id, content, embedding, source, company_id) VALUES ($1, $2, $3, $4, $5)",
		knowledgeBaseId, about.About, about.Embedding, "about", about.CompanyID)

	if err != nil {
		return nil, err
	}
	return &abt, nil
}

func (db *NewDB) GetAbout(id string) (*data.About, error) {
	var abt data.About
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := db.Db.QueryRowContext(ctx, "SELECT id, about FROM about WHERE id = $1", id).Scan(&abt.ID, &abt.About)

	if err != nil {
		return nil, err
	}
	return &abt, nil
}

func (db *NewDB) GetMostRelevantKnowledge(companyID string, queryEmbedding string) (string, error) {
	var content string

	err := db.Db.QueryRow("SELECT content FROM knowledge_base WHERE company_id = $1 ORDER BY embedding <-> $2 LIMIT 1",
		companyID, queryEmbedding).Scan(&content)

	if err != nil {
		return "", err
	}

	return content, nil

}
