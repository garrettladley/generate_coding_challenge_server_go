package storage

import (
	"database/sql"
	"fmt"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/jmoiron/sqlx"
)

type AdminStorage struct {
	Conn *sqlx.DB
}

func NewAdminStorage(conn *sqlx.DB) *AdminStorage {
	return &AdminStorage{Conn: conn}
}

type ApplicantDB struct {
	NUID             sql.NullString `db:"nuid"`
	ApplicantName    sql.NullString `db:"applicant_name"`
	Correct          sql.NullBool   `db:"correct"`
	SubmissionTime   sql.NullTime   `db:"submission_time"`
	RegistrationTime sql.NullTime   `db:"registration_time"`
}

func (s *AdminStorage) Applicant(nuid domain.NUID) (ApplicantDB, error) {
	var applicant ApplicantDB
	err := s.Conn.Get(&applicant, `
	SELECT a.nuid, a.applicant_name, s.correct, s.submission_time, a.registration_time
	FROM applicants a
	LEFT JOIN (
		SELECT nuid, correct, submission_time,
			   ROW_NUMBER() OVER (PARTITION BY nuid ORDER BY submission_time DESC) AS row_num
		FROM submissions
	) s ON a.nuid = s.nuid AND s.row_num = 1
	WHERE a.nuid = $1;
`, nuid)

	if err != nil {
		return ApplicantDB{}, fmt.Errorf("failed to query database: %v", err)
	}

	return applicant, nil
}
