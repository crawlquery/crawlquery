package router_test

import (
	"bytes"
	"crawlquery/api/domain"
	"crawlquery/api/factory"
	"crawlquery/api/router"
	"crawlquery/pkg/authutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthHandler is a mock type for the AuthHandler
type MockAuthHandler struct {
	mock.Mock
}

// Login mocks the Login method of AuthHandler
func (m *MockAuthHandler) Login(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

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

type MockNodeHandler struct {
	mock.Mock
}

func (m *MockNodeHandler) Create(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusCreated, gin.H{"message": "Node created"})
}

func setupRouterWithMocks() map[string]interface{} {
	gin.SetMode(gin.TestMode)

	mockAuthHandler := new(MockAuthHandler)
	mockAuthHandler.On("Login", mock.Anything).Return()

	mockAccountHandler := new(MockAccountHandler)
	mockAccountHandler.On("Create", mock.Anything).Return()

	mockCrawlJobHandler := new(MockCrawlJobHandler)
	mockCrawlJobHandler.On("Create", mock.Anything).Return()

	mockNodeHandler := new(MockNodeHandler)
	mockNodeHandler.On("Create", mock.Anything).Return()

	accountService, accountRepo := factory.AccountServiceWithAccount(&domain.Account{})

	// Setup the router with the mock handler
	testRouter := router.NewRouter(
		accountService,
		mockAuthHandler,
		mockAccountHandler,
		mockCrawlJobHandler,
		mockNodeHandler,
	)

	return map[string]interface{}{
		"testRouter":          testRouter,
		"mockAccountHandler":  mockAccountHandler,
		"mockCrawlJobHandler": mockCrawlJobHandler,
		"mockNodeHandler":     mockNodeHandler,
		"accountService":      accountService,
		"accountRepo":         accountRepo,
	}
}

func TestLoginEndpoint(t *testing.T) {

	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new request to the endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(`{"email":"test@example.com","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	testRouter.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
}

func TestAccountCreationEndpoint(t *testing.T) {

	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
	mockAccountHandler := ifs["mockAccountHandler"].(*MockAccountHandler)

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
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
	mockCrawlJobHandler := ifs["mockCrawlJobHandler"].(*MockCrawlJobHandler)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/crawl-jobs", bytes.NewBufferString(`{"url":"http://example.com"}`))
	req.Header.Set("Content-Type", "application/json")

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Crawl job created")

	mockCrawlJobHandler.AssertExpectations(t)
}

func TestNodeCreationEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
	mockNodeHandler := ifs["mockNodeHandler"].(*MockNodeHandler)
	accountRepo := ifs["accountRepo"].(domain.AccountRepository)

	account, err := accountRepo.GetByEmail("test@example.com")

	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}

	token, err := authutil.GenerateToken(account.ID)

	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/nodes", bytes.NewBufferString(`{"url":"http://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Node created")

	mockNodeHandler.AssertExpectations(t)
}
