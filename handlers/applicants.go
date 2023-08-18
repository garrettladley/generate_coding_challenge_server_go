package handlers

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
)

type ApplicantHandler struct {
	Storage storage.ApplicantStorage
}

func NewApplicantHandler(storage *storage.ApplicantStorage) *ApplicantHandler {
	return &ApplicantHandler{Storage: *storage}
}

func (a *ApplicantHandler) Register(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (a *ApplicantHandler) ForgotToken(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (a *ApplicantHandler) Challenge(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (a *ApplicantHandler) Submit(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
