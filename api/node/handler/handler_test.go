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

func TestCreate(t *testing.T) {
	t.Run("should create a node", func(t *testing.T) {

		account := &domain.Account{
			ID: util.UUID(),
		}
		accSvc, _ := factory.AccountServiceWithAccount(account)

		shardRepo := shardRepo.NewRepository()
		shardRepo.Create(&domain.Shard{
			ID: 0,
		})
		shardSvc := shardSvc.NewService(shardRepo, testutil.NewTestLogger())

		repo := mem.NewRepository()
		svc := service.NewService(repo, accSvc, shardSvc, testutil.NewTestLogger())
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

		account := &domain.Account{
			ID: util.UUID(),
		}
		accSvc, _ := factory.AccountServiceWithAccount(account)

		repo := mem.NewRepository()
		svc := service.NewService(repo, accSvc, nil, testutil.NewTestLogger())
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

		account := &domain.Account{
			ID: util.UUID(),
		}
		accSvc, _ := factory.AccountServiceWithAccount(account)

		repo := mem.NewRepository()
		svc := service.NewService(repo, accSvc, nil, testutil.NewTestLogger())
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
