package integrationtests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"
)

func TestHealthCheckWorks(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/health_check", app.Address), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code 200, but got: %v", resp.StatusCode)
	}
}
