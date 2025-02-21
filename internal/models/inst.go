package models

import "database/sql"

type AppModel struct {
	Model *NewDB
}

func NewAppModel(db *sql.DB) *AppModel {
	return &AppModel{Model: NewDBConn(db)}
}
