package router_test

import (
	"bytes"
	"crawlquery/api/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountHandler is a mock type for the AccountHandler
type MockAccountHandler struct {
	mock.Mock
}

// Create mocks the Create method of AccountHandler
func (m *MockAccountHandler) Create(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Account created"})
}

func TestAccountCreationEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock handler
	mockAccountHandler := new(MockAccountHandler)
	mockAccountHandler.On("Create", mock.Anything).Return()

	// Setup the router with the mock handler
	testRouter := router.NewRouter(mockAccountHandler)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new request to the endpoint
	req, _ := http.NewRequest("POST", "/account", bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	testRouter.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Account created")

	// Check that the mock was called
	mockAccountHandler.AssertExpectations(t)
}
