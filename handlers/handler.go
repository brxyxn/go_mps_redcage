package handlers

import (
	"database/sql"
	"log"
)

type Handlers struct {
	db *sql.DB
	l  *log.Logger
}

func NewHandlers(db *sql.DB, l *log.Logger) *Handlers {
	return &Handlers{db, l}
}
