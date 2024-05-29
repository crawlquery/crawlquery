package handler_test

import (
	"bytes"
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/factory"
	"crawlquery/api/node/handler"
	"crawlquery/api/node/repository/mem"
	"crawlquery/api/node/service"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardSvc "crawlquery/api/shard/service"

	"github.com/gin-gonic/gin"
)

func setup() (*mem.Repository, *service.Service, *domain.Account) {
	account := &domain.Account{
		ID: util.UUIDString(),
	}
	accSvc, _ := factory.AccountServiceWithAccount(account)

	shardRepo := shardRepo.NewRepository()
	shardRepo.Create(&domain.Shard{
		ID: 0,
	})
	shardSvc := shardSvc.NewService(
		shardSvc.WithRepo(shardRepo),
		shardSvc.WithLogger(testutil.NewTestLogger()),
	)
	repo := mem.NewRepository()
	nodeService := service.NewService(
		service.WithAccountService(accSvc),
		service.WithShardService(shardSvc),
		service.WithNodeRepo(repo),
		service.WithLogger(testutil.NewTestLogger()),
	)

	return repo, nodeService, account
}

func TestCreate(t *testing.T) {
	t.Run("should create a node", func(t *testing.T) {

		_, svc, account := setup()
		handler := handler.NewHandler(svc)

		// given
		a := &dto.CreateNodeRequest{
			Hostname: "localhost",
			Port:     8080,
		}

		req, err := a.ToJSON()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Set("account", account)
		ctx.Request = httptest.NewRequest("POST", "/nodes", bytes.NewBuffer(req))

		// when
		handler.Create(ctx)

		// then
		if ctx.Writer.Status() != http.StatusCreated {
			t.Errorf("Expected status to be 201, got %d", ctx.Writer.Status())
		}

		var res dto.CreateNodeResponse
		err = json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Node.Hostname != a.Hostname {
			t.Errorf("Expected hostname to be %s, got %s", a.Hostname, res.Node.Hostname)
		}

		if res.Node.Port != a.Port {
			t.Errorf("Expected port to be %d, got %d", a.Port, res.Node.Port)
		}
	})

	t.Run("should return 400 if malformed JSON", func(t *testing.T) {

		_, svc, account := setup()
		handler := handler.NewHandler(svc)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Set("account", account)
		ctx.Request = httptest.NewRequest("POST", "/nodes", bytes.NewBuffer([]byte("{")))

		handler.Create(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})

	t.Run("should return 400 if hostname is empty", func(t *testing.T) {

		_, svc, account := setup()
		handler := handler.NewHandler(svc)

		a := &dto.CreateNodeRequest{
			Hostname: "",
			Port:     8080,
		}

		req, err := a.ToJSON()

		if err != nil {
			t.Fatalf("Error converting to JSON: %v", err)
		}

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Set("account", account)
		ctx.Request = httptest.NewRequest("POST", "/nodes", bytes.NewBuffer(req))

		handler.Create(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})
}

func TestAuth(t *testing.T) {
	t.Run("should return 401 if no account", func(t *testing.T) {

		_, svc, _ := setup()
		handler := handler.NewHandler(svc)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		body := `{"key":"123"}`
		ctx.Request = httptest.NewRequest("POST", "/auth/node", bytes.NewBuffer([]byte(body)))

		handler.Auth(ctx)

		if ctx.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status to be 401, got %d", ctx.Writer.Status())
		}
	})

	t.Run("should return 400 if no key set", func(t *testing.T) {

		_, svc, account := setup()
		handler := handler.NewHandler(svc)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Set("account", account)
		body := `{}`
		ctx.Request = httptest.NewRequest("POST", "/auth/node", bytes.NewBuffer([]byte(body)))

		handler.Auth(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("Expected status to be 400, got %d", ctx.Writer.Status())
		}
	})

	t.Run("should return node if key is correct", func(t *testing.T) {

		repo, svc, account := setup()
		handler := handler.NewHandler(svc)

		node := &domain.Node{
			ID:        util.UUIDString(),
			Key:       "123",
			AccountID: account.ID,
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   1,
		}

		repo.Create(node)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Set("account", account)
		body := `{"key":"123"}`
		ctx.Request = httptest.NewRequest("POST", "/auth/node", bytes.NewBuffer([]byte(body)))

		handler.Auth(ctx)

		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status to be 200, got %d", ctx.Writer.Status())
		}

		var res dto.AuthenticateNodeResponse
		err := json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if res.Node.ID != node.ID {
			t.Errorf("Expected ID to be %s, got %s", node.ID, res.Node.ID)
		}

		if res.Node.Key != node.Key {
			t.Errorf("Expected Key to be %s, got %s", node.Key, res.Node.Key)
		}

		if res.Node.AccountID != node.AccountID {
			t.Errorf("Expected AccountID to be %s, got %s", node.AccountID, res.Node.AccountID)
		}

		if res.Node.Hostname != node.Hostname {
			t.Errorf("Expected Hostname to be %s, got %s", node.Hostname, res.Node.Hostname)
		}

		if res.Node.Port != node.Port {
			t.Errorf("Expected Port to be %d, got %d", node.Port, res.Node.Port)
		}

		if domain.ShardID(res.Node.ShardID) != node.ShardID {
			t.Errorf("Expected ShardID to be %d, got %d", node.ShardID, res.Node.ShardID)
		}
	})
}

func TestListByShard(t *testing.T) {
	t.Run("should return nodes by shard", func(t *testing.T) {

		repo, svc, account := setup()
		handler := handler.NewHandler(svc)

		node1 := &domain.Node{
			ID:        util.UUIDString(),
			Key:       "123",
			AccountID: account.ID,
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   0,
		}

		node2 := &domain.Node{
			ID:        util.UUIDString(),
			Key:       "123",
			AccountID: account.ID,
			Hostname:  "localhost",
			Port:      8080,
			ShardID:   0,
		}

		repo.Create(node1)
		repo.Create(node2)

		responseWriter := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(responseWriter)
		ctx.Params = []gin.Param{
			{
				Key:   "shardID",
				Value: "0",
			},
		}

		handler.ListByShardID(ctx)

		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status to be 200, got %d", ctx.Writer.Status())
		}

		var res dto.ListNodesByShardResponse
		err := json.NewDecoder(responseWriter.Body).Decode(&res)

		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if len(res.Nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(res.Nodes))
		}
	})
}
