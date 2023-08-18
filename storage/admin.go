package storage

import (
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

type ApplicantFound struct {
	Applicant        domain.Applicant
	Correct          bool
	TimeToCompletion time.Duration
}

type GetApplicantsResult struct {
	ApplicantsFound        []ApplicantFound
	ApplicantsNotSubmitted []domain.NUID
	ApplicantsNotFound     []domain.NUID
}

func (s *AdminStorage) GetApplicants(nuids []domain.NUID) (GetApplicantsResult, error) {
	applicants := []ApplicantFound{}
	err := s.Conn.Select(applicants, "SELECT DISTINCT ON (nuid) nuid, applicant_name, correct, submission_time, registration_time FROM submissions JOIN applicants using(nuid) where nuid=ANY(?) ORDER BY nuid, submission_time DESC;", nuids)

	if err != nil {
		return GetApplicantsResult{}, err
	}

	if len(nuids) != len(applicants) {
		return notAllApplicantsSubmittedHandler(s, nuids, applicants)
	}
	return GetApplicantsResult{
		ApplicantsFound:        applicants,
		ApplicantsNotSubmitted: []domain.NUID{},
		ApplicantsNotFound:     []domain.NUID{},
	}, nil
}

func notAllApplicantsSubmittedHandler(s *AdminStorage, nuids []domain.NUID, applicants []ApplicantFound) (GetApplicantsResult, error) {
	applicantNuids := []domain.NUID{}
	for _, applicant := range applicants {
		applicantNuids = append(applicantNuids, applicant.Applicant.NUID)
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
