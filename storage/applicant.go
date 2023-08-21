package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

type RegisterResult struct {
	Token      *uuid.UUID `json:"token,omitempty"`
	Challenge  []string   `json:"challenge,omitempty"`
	Message    string     `json:"message,omitempty"`
	HttpStatus int        `json:"-"`
}

func (s *ApplicantStorage) Register(applicant domain.Applicant) (RegisterResult, error) {
	registrationTime := time.Now()
	token := uuid.New()
	red, _ := domain.Red.String()
	orange, _ := domain.Orange.String()
	yellow, _ := domain.Yellow.String()
	green, _ := domain.Green.String()
	blue, _ := domain.Blue.String()
	violet, _ := domain.Violet.String()
	challenge := domain.GenerateChallenge(100, []string{"", red, orange, yellow, green, blue, violet})

	insertSataement := "INSERT INTO applicants (nuid, applicant_name, registration_time, token, challenge, solution) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := s.Conn.Exec(insertSataement, applicant.NUID, applicant.Name, registrationTime, token, pq.Array(challenge.Challenge), pq.Array(challenge.Solution))

	if err != nil {
		pgErr, isPGError := err.(*pq.Error)
		if isPGError && pgErr.Code == "23505" {
			return RegisterResult{Message: fmt.Sprintf("NUID %s has already registered! Use the forgot_token endpoint to retrieve your token.", applicant.NUID), HttpStatus: 409}, nil
		}
		return RegisterResult{}, err
	}

	return RegisterResult{Token: &token, Challenge: challenge.Challenge}, nil
}

type ForgotTokenResult struct {
	Message    string    `json:"message,omitempty"`
	Token      uuid.UUID `json:"token"`
	HttpStatus int       `json:"-"`
}

type ForgotTokenDB struct {
	Token sql.NullString `db:"token"`
}

func (s *ApplicantStorage) ForgotToken(nuid domain.NUID) (ForgotTokenResult, error) {
	var dbResult ForgotTokenDB
	err := s.Conn.Get(&dbResult, "SELECT token FROM applicants WHERE nuid=$1;", nuid)

	if !dbResult.Token.Valid {
		return ForgotTokenResult{Message: fmt.Sprintf("Applicant with NUID %s not found!", nuid), HttpStatus: 404}, nil
	} else if err != nil {
		return ForgotTokenResult{}, err
	}

	token, err := uuid.Parse(dbResult.Token.String)

	if err != nil {
		return ForgotTokenResult{}, err
	}

	return ForgotTokenResult{Token: token}, nil
}

type ChallengeResult struct {
	Message    string   `json:"message,omitempty"`
	Challenge  []string `json:"challenge"`
	HttpStatus int      `json:"-"`
}

type ChallengeDB struct {
	Challenge StringArray `db:"challenge"`
}

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	strArray := strings.Split(string(src.([]byte)), ",")
	for i, s := range strArray {
		s = strings.Trim(s, "\"")

		if s == "\\\"" {
			s = ""
		}

		strArray[i] = s
	}

	first := strArray[0]
	if strings.HasPrefix(first, "{") {
		strArray[0] = strings.TrimPrefix(first, "{")
	} else {
		return errors.New("invalid array format: first element does not start with '{'")
	}

	lastIndex := len(strArray) - 1
	last := strArray[lastIndex]
	if strings.HasSuffix(last, "}") {
		strArray[lastIndex] = strings.TrimSuffix(last, "}")
	} else {
		return errors.New("invalid array format: last element does not end with '}'")
	}

	*a = strArray
	return nil
}

func (s *ApplicantStorage) Challenge(token uuid.UUID) (ChallengeResult, error) {
	var dbResult ChallengeDB
	err := s.Conn.Get(&dbResult, "SELECT challenge FROM applicants WHERE token=$1;", token)

	if len(dbResult.Challenge) == 0 && err != nil {
		return ChallengeResult{Message: fmt.Sprintf("Record associated with token %s not found!", token), HttpStatus: 404}, nil
	} else if err != nil {
		return ChallengeResult{}, err
	}

	return ChallengeResult{Challenge: dbResult.Challenge}, nil
}

type SubmitResult struct {
	Message    string `json:"message,omitempty"`
	Correct    bool   `json:"correct,omitempty"`
	HttpStatus int    `json:"-"`
}
type SubmitDB struct {
	NUID     sql.NullString `db:"nuid"`
	Solution StringArray    `db:"solution"`
}

type WriteSubmitDB struct {
	NUID    domain.NUID
	Correct bool
}

func (s *ApplicantStorage) Submit(token uuid.UUID, givenSolution []string) (SubmitResult, error) {
	var dbResult SubmitDB
	err := s.Conn.Get(&dbResult, "SELECT nuid, solution FROM applicants WHERE token=$1;", token)
	if err != nil {
		return SubmitResult{}, err
	}

	if !dbResult.NUID.Valid && len(dbResult.Solution) == 0 {
		return SubmitResult{Message: fmt.Sprintf("Record associated with token %s not found!", token), HttpStatus: 404}, nil
	}

	nuid, err := domain.ParseNUID(dbResult.NUID.String)

	if err != nil {
		return SubmitResult{}, fmt.Errorf("invalid database state!. Error: %v", err)
	}

	correct, err := s.writeSubmit(*nuid, givenSolution, dbResult.Solution)

	if err != nil {
		return SubmitResult{}, err
	}

	return SubmitResult{Correct: correct}, nil
}

func (s *ApplicantStorage) writeSubmit(nuid domain.NUID, givenSolution []string, actualSolution []string) (bool, error) {
	var correct bool
	if len(givenSolution) == len(actualSolution) {
		correct = true
		for i := range givenSolution {
			if givenSolution[i] != actualSolution[i] {
				correct = false
				break
			}
		}
	} else {
		correct = false
	}

	insertStatement := "INSERT INTO submissions (nuid, correct, submission_time) VALUES ($1, $2, $3);"
	_, err := s.Conn.Exec(insertStatement, nuid, correct, time.Now())

	if err != nil {
		return correct, err
	}

	return correct, err
}
