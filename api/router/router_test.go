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

type MockPageHandler struct {
	mock.Mock
}

func (m *MockPageHandler) Create(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusCreated, gin.H{"message": "Page created"})
}

type MockNodeHandler struct {
	mock.Mock
}

func (m *MockNodeHandler) Create(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusCreated, gin.H{"message": "Node created"})
}

func (m *MockNodeHandler) ListByAccountID(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Nodes listed"})
}

func (m *MockNodeHandler) Auth(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Node authenticated"})
}

func (m *MockNodeHandler) ListByShardID(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Nodes listed"})
}

type MockSearchHandler struct {
	mock.Mock
}

func (m *MockSearchHandler) Search(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Search successful"})
}

func setupRouterWithMocks() map[string]interface{} {
	gin.SetMode(gin.TestMode)

	mockAuthHandler := new(MockAuthHandler)
	mockAuthHandler.On("Login", mock.Anything).Return()

	mockAccountHandler := new(MockAccountHandler)
	mockAccountHandler.On("Create", mock.Anything).Return()

	mockPageHandler := new(MockPageHandler)
	mockPageHandler.On("Create", mock.Anything).Return()

	mockNodeHandler := new(MockNodeHandler)
	mockNodeHandler.On("Create", mock.Anything).Return()
	mockNodeHandler.On("ListByAccountID", mock.Anything).Return()
	mockNodeHandler.On("ListByShardID", mock.Anything).Return()
	mockNodeHandler.On("Auth", mock.Anything).Return()

	mockSearchHandler := new(MockSearchHandler)
	mockSearchHandler.On("Search", mock.Anything).Return()

	accountService, accountRepo := factory.AccountServiceWithAccount(&domain.Account{})

	// Setup the router with the mock handler
	testRouter := router.NewRouter(
		accountService,
		mockAuthHandler,
		mockAccountHandler,
		mockPageHandler,
		mockNodeHandler,
		mockSearchHandler,
	)

	return map[string]interface{}{
		"testRouter":         testRouter,
		"mockAccountHandler": mockAccountHandler,
		"mockPageHandler":    mockPageHandler,
		"mockNodeHandler":    mockNodeHandler,
		"mockSearchHandler":  mockSearchHandler,
		"accountService":     accountService,
		"accountRepo":        accountRepo,
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

func TestCreatePageEndoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
	mockPageHandler := ifs["mockPageHandler"].(*MockPageHandler)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/pages", bytes.NewBufferString(`{"url":"http://example.com"}`))
	req.Header.Set("Content-Type", "application/json")

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Page created")

	mockPageHandler.AssertExpectations(t)
}

func TestNodeCreationEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
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
}

func TestSearchEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
	mockSearchHandler := ifs["mockSearchHandler"].(*MockSearchHandler)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/search?q=term", nil)

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Search successful")

	mockSearchHandler.AssertExpectations(t)
}

func TestNodeListByAccountIDEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)
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

	req, _ := http.NewRequest("GET", "/nodes", nil)

	req.Header.Set("Authorization", "Bearer "+token)

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "Nodes listed")
}

func TestNodeAuthEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/auth/node", nil)

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNodeListByShardIDEndpoint(t *testing.T) {
	// Set the router to test mode
	ifs := setupRouterWithMocks()

	testRouter := ifs["testRouter"].(*gin.Engine)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/shards/1/nodes", nil)

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Nodes listed")
}
