package handlers

import (
	"fmt"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
)

type ApplicantHandler storage.ApplicantStorage

func NewApplicantHandler(storage *storage.ApplicantStorage) *ApplicantHandler {
	return (*ApplicantHandler)(storage)
}

type RegisterRequestBody struct {
	RawNUID          string `json:"nuid"`
	RawApplicantName string `json:"name"`
}

func (a *ApplicantHandler) Register(c *fiber.Ctx) error {
	var registerRequestBody RegisterRequestBody

	if err := c.BodyParser(&registerRequestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body %s", registerRequestBody))
	}

	nuid, err := domain.ParseNUID(registerRequestBody.RawNUID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid nuid %s", registerRequestBody.RawNUID))
	}

	applicantName, err := domain.ParseApplicantName(registerRequestBody.RawApplicantName)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid applicant name %s", registerRequestBody.RawApplicantName))
	}

	result, err := (*storage.ApplicantStorage)(a).RegisterApplicant(domain.Applicant{
		NUID: *nuid,
		Name: *applicantName,
	})

	if err != nil {
		return err
	}

	if result.HttpStatus == 409 {
		return c.Status(result.HttpStatus).JSON(result)
	}

	return c.JSON(result)
}

func (a *ApplicantHandler) ForgotToken(c *fiber.Ctx) error {
	rawNUID := c.Params("nuid")

	nuid, err := domain.ParseNUID(rawNUID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid nuid %s", rawNUID))
	}

	result, err := (*storage.ApplicantStorage)(a).ForgotToken(*nuid)

	if err != nil {
		return err
	}

	if result.HttpStatus == 404 {
		return c.Status(result.HttpStatus).JSON("Record associated with given NUID not found!")
	}

	return c.JSON(result)
}

func (a *ApplicantHandler) Challenge(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (a *ApplicantHandler) Submit(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
