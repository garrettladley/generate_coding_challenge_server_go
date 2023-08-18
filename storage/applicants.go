package storage

import (
	"github.com/jmoiron/sqlx"
)

type ApplicantStorage struct {
	Conn *sqlx.DB
}

func NewApplicantStorage(conn *sqlx.DB) *ApplicantStorage {
	return &ApplicantStorage{Conn: conn}
}
