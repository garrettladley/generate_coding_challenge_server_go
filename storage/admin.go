package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/jmoiron/sqlx"
)

type AdminStorage struct {
	Conn *sqlx.DB
}

func NewAdminStorage(conn *sqlx.DB) *AdminStorage {
	return &AdminStorage{Conn: conn}
}

type GetApplicantDB struct {
	NUID             sql.NullString `db:"nuid"`
	ApplicantName    sql.NullString `db:"applicant_name"`
	Correct          sql.NullBool   `db:"correct"`
	SubmissionTime   sql.NullTime   `db:"submission_time"`
	RegistrationTime sql.NullTime   `db:"registration_time"`
}

type ApplicantFound struct {
	NUID             domain.NUID
	ApplicantName    domain.ApplicantName
	Correct          bool
	TimeToCompletion time.Duration
}

type GetApplicantResult struct {
	Message          string               `json:"message,omitempty"`
	NUID             domain.NUID          `json:"nuid,omitempty"`
	Name             domain.ApplicantName `json:"name,omitempty"`
	Correct          *bool                `json:"correct,omitempty"`
	TimeToCompletion TimeToCompletion     `json:"time_to_completion,omitempty"`
	HttpStatus       int                  `json:"-"`
}

type TimeToCompletion struct {
	Seconds int `json:"seconds"`
	Nanos   int `json:"nanos"`
}

func convert(t time.Duration) TimeToCompletion {
	return TimeToCompletion{
		Seconds: int(t.Seconds()),
		Nanos:   int(t.Nanoseconds()),
	}
}

func (s *AdminStorage) Applicant(nuid domain.NUID) (GetApplicantResult, error) {
	var applicant GetApplicantDB
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

	if !applicant.NUID.Valid && !applicant.ApplicantName.Valid && !applicant.Correct.Valid && !applicant.SubmissionTime.Valid && !applicant.RegistrationTime.Valid && err != nil {
		return GetApplicantResult{Message: fmt.Sprintf("Applicant with NUID %s not found!", nuid), HttpStatus: 404}, nil
	} else if err != nil {
		return GetApplicantResult{}, err
	}

	if !applicant.Correct.Valid && !applicant.SubmissionTime.Valid {
		return GetApplicantResult{Message: fmt.Sprintf("Applicant with NUID %s has not submitted yet!", nuid)}, nil
	}

	applicantFound, err := processGetApplicantDB(applicant)
	if err != nil {
		return GetApplicantResult{}, fmt.Errorf("invalid database state!. Error: %v", err)
	}

	return GetApplicantResult{
		NUID:             applicantFound.NUID,
		Name:             applicantFound.ApplicantName,
		Correct:          &applicantFound.Correct,
		TimeToCompletion: convert(applicantFound.TimeToCompletion),
	}, nil
}

func processGetApplicantDB(applicant GetApplicantDB) (ApplicantFound, error) {
	nuid, err := domain.ParseNUID(applicant.NUID.String)

	if err != nil {
		return ApplicantFound{}, err
	}

	applicantName, err := domain.ParseApplicantName(applicant.ApplicantName.String)

	if err != nil {
		return ApplicantFound{}, err
	}

	return ApplicantFound{
		NUID:             *nuid,
		ApplicantName:    *applicantName,
		Correct:          applicant.Correct.Bool,
		TimeToCompletion: applicant.SubmissionTime.Time.Sub(applicant.RegistrationTime.Time),
	}, nil
}
