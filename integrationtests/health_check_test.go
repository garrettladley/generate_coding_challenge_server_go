package integrationtests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckWorks(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/health_check", app.Address), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	assert.Equal(200, resp.StatusCode)
}
