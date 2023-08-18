package handlers

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler storage.AdminStorage

func NewAdminHandler(storage *storage.AdminStorage) *AdminHandler {
	return (*AdminHandler)(storage)
}

// Applicants godoc
// @Summary Get the status of the provided applicants
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} storage.GetApplicantsResult
// @Router /admin/applicants [get]
func (u *AdminHandler) Applicants(c *fiber.Ctx) error {
	var rawNUIDs []string

	err := c.BodyParser(&rawNUIDs)
	if err != nil {
		return err
	}

	nuids := make([]domain.NUID, len(rawNUIDs))

	for idx, rawNUID := range rawNUIDs {
		nuid, err := domain.ParseNUID(rawNUID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid nuid %s", rawNUID)
		}
		nuids[idx] = *nuid
	}

	result, err := (*storage.AdminStorage)(u).GetApplicants(nuids)

	if err != nil {
		return err
	}

	if len(result.ApplicantsNotFound) != 0 {
		return c.Status(fiber.StatusNotFound).JSON(result)
	} else {
		return c.JSON(result)
	}
}
