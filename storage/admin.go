package storage

import (
	"fmt"
	"time"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AdminStorage struct {
	Conn *sqlx.DB
}

func NewAdminStorage(conn *sqlx.DB) *AdminStorage {
	return &AdminStorage{Conn: conn}
}

type ApplicantFoundDB struct {
	NUID             string    `db:"nuid"`
	ApplicantName    string    `db:"applicant_name"`
	Correct          bool      `db:"correct"`
	SubmissionTime   time.Time `db:"submission_time"`
	RegistrationTime time.Time `db:"registration_time"`
}

type ApplicantFound struct {
	NUID             domain.NUID
	ApplicantName    domain.ApplicantName
	Correct          bool
	TimeToCompletion time.Duration
}

type GetApplicantsResult struct {
	ApplicantsFound        []ApplicantFound
	ApplicantsNotSubmitted []domain.NUID
	ApplicantsNotFound     []domain.NUID
}

func (s *AdminStorage) GetApplicants(nuids []domain.NUID) (GetApplicantsResult, error) {
	applicantsFoundDB := []ApplicantFoundDB{}
	err := s.Conn.Select(&applicantsFoundDB, `
		SELECT DISTINCT ON (nuid) nuid, applicant_name, correct, submission_time, registration_time 
		FROM submissions 
		JOIN applicants USING (nuid) 
		WHERE nuid = ANY($1) 
		ORDER BY nuid, submission_time DESC;
	`, pq.Array(nuids))

	if err != nil {
		return GetApplicantsResult{}, err
	}

	applicantsFound := []ApplicantFound{}
	for _, applicant := range applicantsFoundDB {
		applicantFound, err := processApplicantFoundDB(applicant)
		if err != nil {
			return GetApplicantsResult{}, fmt.Errorf("invalid database state!. Error: %v", err)
		}
		applicantsFound = append(applicantsFound, applicantFound)
	}

	if len(nuids) != len(applicantsFound) {
		return notAllApplicantsSubmittedHandler(s, nuids, applicantsFound)
	}

	return GetApplicantsResult{
		ApplicantsFound:        applicantsFound,
		ApplicantsNotSubmitted: []domain.NUID{},
		ApplicantsNotFound:     []domain.NUID{},
	}, nil
}

func processApplicantFoundDB(applicant ApplicantFoundDB) (ApplicantFound, error) {
	nuid, err := domain.ParseNUID(applicant.NUID)

	if err != nil {
		return ApplicantFound{}, err
	}

	applicantName, err := domain.ParseApplicantName(applicant.ApplicantName)

	if err != nil {
		return ApplicantFound{}, err
	}

	return ApplicantFound{
		NUID:             *nuid,
		ApplicantName:    *applicantName,
		Correct:          applicant.Correct,
		TimeToCompletion: applicant.SubmissionTime.Sub(applicant.RegistrationTime),
	}, nil
}

func notAllApplicantsSubmittedHandler(s *AdminStorage, nuids []domain.NUID, applicants []ApplicantFound) (GetApplicantsResult, error) {
	applicantNuids := []domain.NUID{}
	for _, applicant := range applicants {
		applicantNuids = append(applicantNuids, applicant.NUID)
	}
	remainingNuids := findElementsNotInB(nuids, applicantNuids)

	applicantsNotSubmitted := []domain.NUID{}
	applicantsNotFound := []domain.NUID{}

	for _, nuid := range remainingNuids {
		_, err := s.GetApplicant(nuid)
		if err != nil {
			applicantsNotFound = append(applicantsNotFound, nuid)
		} else {
			applicantsNotSubmitted = append(applicantsNotSubmitted, nuid)
		}
	}

	return GetApplicantsResult{
		ApplicantsFound:        applicants,
		ApplicantsNotSubmitted: applicantsNotSubmitted,
		ApplicantsNotFound:     applicantsNotFound,
	}, nil
}

func findElementsNotInB[T comparable](listA, listB []T) []T {
	notInB := []T{}

	for _, elementA := range listA {
		found := false
		for _, elementB := range listB {
			if elementA == elementB {
				found = true
				break
			}
		}

		if !found {
			notInB = append(notInB, elementA)
		}
	}

	return notInB
}

func (s *AdminStorage) GetApplicant(nuid domain.NUID) (ApplicantFound, error) {
	applicant := ApplicantFound{}
	err := s.Conn.Get(applicant, "SELECT nuid FROM applicants where nuid=$1;", nuid)

	if err != nil {
		return applicant, err
	}
	return applicant, nil
}
