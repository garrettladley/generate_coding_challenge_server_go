package handlers

import (
	"fmt"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ApplicantHandler storage.ApplicantStorage

func NewApplicantHandler(storage *storage.ApplicantStorage) *ApplicantHandler {
	return (*ApplicantHandler)(storage)
}

type RegisterRequestBody struct {
	RawApplicantName string `json:"name"`
	RawNUID          string `json:"nuid"`
}

type RegisterResponse struct {
	Token     uuid.UUID `json:"token"`
	Challenge []string  `json:"challenge"`
}

func (a *ApplicantHandler) Register(c *fiber.Ctx) error {
	var registerRequestBody RegisterRequestBody

	if err := c.BodyParser(&registerRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid request body %s", registerRequestBody))
	}

	nuid, err := domain.ParseNUID(registerRequestBody.RawNUID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid NUID %s", registerRequestBody.RawNUID))
	}

	applicantName, err := domain.ParseApplicantName(registerRequestBody.RawApplicantName)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid applicant name %s", registerRequestBody.RawApplicantName))
	}

	result, err := (*storage.ApplicantStorage)(a).Register(domain.Applicant{
		NUID: *nuid,
		Name: *applicantName,
	})

	if err != nil {
		pgErr, isPGError := err.(*pq.Error)
		if isPGError && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).SendString(fmt.Sprintf("NUID %s has already registered! Use the forgot_token endpoint to retrieve your token.", nuid))
		} else {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

type ForgotTokenResponse struct {
	Token uuid.UUID `json:"token"`
}

func (a *ApplicantHandler) ForgotToken(c *fiber.Ctx) error {
	rawNUID := c.Params("nuid")

	nuid, err := domain.ParseNUID(rawNUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid NUID %s", rawNUID))
	}

	result, err := (*storage.ApplicantStorage)(a).ForgotToken(*nuid)

	if err != nil && !result.Token.Valid {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Applicant with NUID %s not found!", nuid))
	} else if err != nil {
		return err
	}

	token, err := uuid.Parse(result.Token.String)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("invalid database state! Error: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(ForgotTokenResponse{
		Token: token,
	})
}

type ChallengeResponse struct {
	Challenge []string `json:"challenge"`
}

func (a *ApplicantHandler) Challenge(c *fiber.Ctx) error {
	rawToken := c.Params("token")

	token, err := uuid.Parse(rawToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid token %s", rawToken))
	}

	result, err := (*storage.ApplicantStorage)(a).Challenge(token)

	if err != nil && len(result.Challenge) == 0 {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Record associated with token %s not found!", token))
	} else if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ChallengeResponse{
		Challenge: result.Challenge,
	})
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

	if !result.NUID.Valid && len(result.Solution) == 0 {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Record associated with token %s not found!", token))
	}

	nuid, err := domain.ParseNUID(result.NUID.String)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("invalid database state! Error: %v", err))
	}

	correct, err := (*storage.ApplicantStorage)(a).WriteSubmit(*nuid, result.Solution, submitRequestBody)

	if err != nil {
		return err
	}

	if correct {
		return c.JSON(SubmitResponseBody{
			Correct: correct,
			Message: "Correct - nice work!",
		})
	} else {
		return c.JSON(SubmitResponseBody{
			Correct: correct,
			Message: "Incorrect Solution",
		})
	}
}
