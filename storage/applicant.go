package storage

import (
	"database/sql"
	"errors"
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
	Token     uuid.UUID
	Challenge []string
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
		return RegisterResult{}, err
	}

	return RegisterResult{Token: token, Challenge: challenge.Challenge}, nil
}

type ForgotTokenDB struct {
	Token sql.NullString `db:"token"`
}

func (s *ApplicantStorage) ForgotToken(nuid domain.NUID) (ForgotTokenDB, error) {
	var dbResult ForgotTokenDB
	err := s.Conn.Get(&dbResult, "SELECT token FROM applicants WHERE nuid=$1;", nuid)

	if err != nil {
		return ForgotTokenDB{}, err
	}

	return dbResult, nil
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

func (s *ApplicantStorage) Challenge(token uuid.UUID) (ChallengeDB, error) {
	var dbResult ChallengeDB
	err := s.Conn.Get(&dbResult, "SELECT challenge FROM applicants WHERE token=$1;", token)

	if err != nil {
		return ChallengeDB{}, err
	}

	return dbResult, nil
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

func (s *ApplicantStorage) Submit(token uuid.UUID, givenSolution []string) (SubmitDB, error) {
	var dbResult SubmitDB
	err := s.Conn.Get(&dbResult, "SELECT nuid, solution FROM applicants WHERE token=$1;", token)

	if err != nil {
		return SubmitDB{}, err
	}

	return dbResult, nil
}

func (s *ApplicantStorage) WriteSubmit(nuid domain.NUID, givenSolution []string, actualSolution []string) (bool, error) {
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
