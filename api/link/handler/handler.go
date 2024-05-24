package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	linkService domain.LinkService
	logger      *zap.SugaredLogger
}

func NewHandler(linkService domain.LinkService, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		linkService: linkService,
		logger:      logger,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	_, err := h.linkService.Create(req.Src, req.Dst)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
}