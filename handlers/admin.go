package handlers

import (
	"fmt"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler storage.AdminStorage

func NewAdminHandler(storage *storage.AdminStorage) *AdminHandler {
	return (*AdminHandler)(storage)
}

func (a *AdminHandler) Applicant(c *fiber.Ctx) error {
	rawNUID := c.Params("nuid")

	nuid, err := domain.ParseNUID(rawNUID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid nuid %s", rawNUID))
	}

	result, err := (*storage.AdminStorage)(a).Applicant(*nuid)

	if err != nil {
		return err
	}

	if result.HttpStatus == 404 {
		return c.Status(result.HttpStatus).JSON(result)
	}

	return c.JSON(result)
}
