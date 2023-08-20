package handlers

import (
	"fmt"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	rawToken := c.Params("token")

	token, err := uuid.Parse(rawToken)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid token %s", rawToken))
	}

	result, err := (*storage.ApplicantStorage)(a).Challenge(token)

	if err != nil {
		return err
	}

	if result.HttpStatus == 404 {
		return c.Status(result.HttpStatus).JSON("Record associated with given token not found!")
	}

	return c.JSON(result)
}

type SubmitRequestBody []string

type SubmitResponseBody struct {
	Correct bool   `json:"correct"`
	Message string `json:"message"`
}

func (a *ApplicantHandler) Submit(c *fiber.Ctx) error {
	rawToken := c.Params("token")

	token, err := uuid.Parse(rawToken)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid token %s", rawToken))
	}

	var submitRequestBody SubmitRequestBody

	if err := c.BodyParser(&submitRequestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body %s", submitRequestBody))
	}

	result, err := (*storage.ApplicantStorage)(a).Submit(token, submitRequestBody)

	if err != nil {
		return err
	}

	if result.HttpStatus == 404 {
		return c.Status(result.HttpStatus).JSON("Record associated with given token not found!")
	}

	if result.Correct {
		return c.JSON(SubmitResponseBody{
			Correct: result.Correct,
			Message: "Correct - nice work!",
		})
	} else {
		return c.JSON(SubmitResponseBody{
			Correct: result.Correct,
			Message: "Incorrect Solution",
		})
	}
}
