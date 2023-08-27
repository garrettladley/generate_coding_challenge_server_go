package handlers

import (
	"fmt"
	"time"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler storage.AdminStorage

func NewAdminHandler(storage *storage.AdminStorage) *AdminHandler {
	return (*AdminHandler)(storage)
}

type ApplicantResponse struct {
	NUID             domain.NUID          `json:"nuid"`
	ApplicantName    domain.ApplicantName `json:"name"`
	Correct          bool                 `json:"correct"`
	TimeToCompletion TimeToCompletion     `json:"time_to_completion"`
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

func processApplicantDB(applicant storage.ApplicantDB) (ApplicantResponse, error) {
	nuid, err := domain.ParseNUID(applicant.NUID.String)

	if err != nil {
		return ApplicantResponse{}, err
	}

	applicantName, err := domain.ParseApplicantName(applicant.ApplicantName.String)

	if err != nil {
		return ApplicantResponse{}, err
	}

	return ApplicantResponse{
		NUID:             *nuid,
		ApplicantName:    *applicantName,
		Correct:          applicant.Correct.Bool,
		TimeToCompletion: convert(applicant.SubmissionTime.Time.Sub(applicant.RegistrationTime.Time)),
	}, nil
}

func (a *AdminHandler) Applicant(c *fiber.Ctx) error {
	rawNUID := c.Params("nuid")

	nuid, err := domain.ParseNUID(rawNUID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid NUID %s", rawNUID))
	}

	result, err := (*storage.AdminStorage)(a).Applicant(*nuid)

	if err != nil && !result.NUID.Valid && !result.ApplicantName.Valid && !result.Correct.Valid && !result.SubmissionTime.Valid && !result.RegistrationTime.Valid {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Applicant with NUID %s not found!", nuid))
	} else if err != nil {
		return err
	}

	if !result.Correct.Valid && !result.SubmissionTime.Valid {
		return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Applicant with NUID %s has not submitted yet!", nuid))
	}

	ApplicantResponse, err := processApplicantDB(result)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("invalid database state! Error: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(ApplicantResponse)
}
