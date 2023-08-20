package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ApplicantStorage struct {
	Conn *sqlx.DB
}

func NewApplicantStorage(conn *sqlx.DB) *ApplicantStorage {
	return &ApplicantStorage{Conn: conn}
}

type RegisterApplicantResponse struct {
	Token      uuid.UUID `json:"-"`
	Message    string    `json:"message,omitempty"`
	HttpStatus int       `json:"-"`
}

func (r RegisterApplicantResponse) MarshalJSON() ([]byte, error) {
	if r.Token == uuid.Nil {
		type Alias RegisterApplicantResponse
		return json.Marshal(&struct {
			Alias
			Token string `json:"token,omitempty"`
		}{
			Alias: (Alias)(r),
			Token: "",
		})
	}
	return json.Marshal(&r)
}

func (s *ApplicantStorage) RegisterApplicant(applicant domain.Applicant) (RegisterApplicantResponse, error) {
	registrationTime := time.Now()
	token := uuid.New()
	red, _ := domain.Red.String()
	orange, _ := domain.Orange.String()
	yellow, _ := domain.Yellow.String()
	green, _ := domain.Green.String()
	blue, _ := domain.Blue.String()
	violet, _ := domain.Violet.String()
	challenge := domain.GenerateChallenge(applicant.NUID, 100, []string{red, orange, yellow, green, blue, violet})

	insertSataement := "INSERT INTO applicants (nuid, applicant_name, registration_time, token, challenge, solution) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := s.Conn.Exec(insertSataement, applicant.NUID, applicant.Name, registrationTime, token, pq.Array(challenge.Challenge), pq.Array(challenge.Solution))

	if err != nil {
		pgErr, isPGError := err.(*pq.Error)
		if isPGError && pgErr.Code == "23505" {
			return RegisterApplicantResponse{Message: fmt.Sprintf("NUID %s has already registered! Use the forgot_token endpoint to retrieve your token.", applicant.NUID), HttpStatus: 409}, nil
		}
		return RegisterApplicantResponse{}, err
	}

	return RegisterApplicantResponse{Token: token}, nil
}

type ForgotTokenResponse struct {
	Token      uuid.UUID `json:"token"`
	HttpStatus int       `json:"-"`
}

type ForgotTokenDB struct {
	Token sql.NullString `db:"token"`
}

func (s *ApplicantStorage) ForgotToken(nuid domain.NUID) (ForgotTokenResponse, error) {
	var dbResult ForgotTokenDB
	err := s.Conn.Get(&dbResult, "SELECT token FROM applicants WHERE nuid=$1;", nuid)

	if !dbResult.Token.Valid && err != nil {
		return ForgotTokenResponse{HttpStatus: 404}, nil
	} else if err != nil {
		return ForgotTokenResponse{}, err
	}

	token, err := uuid.Parse(dbResult.Token.String)

	if err != nil {
		return ForgotTokenResponse{}, err
	}

	return ForgotTokenResponse{Token: token}, nil
}
