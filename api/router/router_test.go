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

type MockCrawlJobHandler struct {
	mock.Mock
}

func (m *MockCrawlJobHandler) Create(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusCreated, gin.H{"message": "Crawl job created"})
}

func TestAccountCreationEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock handler
	mockAccountHandler := new(MockAccountHandler)
	mockAccountHandler.On("Create", mock.Anything).Return()

	mockCrawlJobHandler := new(MockCrawlJobHandler)
	mockCrawlJobHandler.On("Create", mock.Anything).Return()

	// Setup the router with the mock handler
	testRouter := router.NewRouter(mockAccountHandler, mockCrawlJobHandler)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new request to the endpoint
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	testRouter.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Account created")

	// Check that the mock was called
	mockAccountHandler.AssertExpectations(t)
}

func TestCrawlJobCreationEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAccountHandler := new(MockAccountHandler)
	mockAccountHandler.On("Create", mock.Anything).Return()

	mockCrawlJobHandler := new(MockCrawlJobHandler)
	mockCrawlJobHandler.On("Create", mock.Anything).Return()

	testRouter := router.NewRouter(mockAccountHandler, mockCrawlJobHandler)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/crawl-jobs", bytes.NewBufferString(`{"url":"http://example.com"}`))
	req.Header.Set("Content-Type", "application/json")

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Crawl job created")

	mockCrawlJobHandler.AssertExpectations(t)
}
