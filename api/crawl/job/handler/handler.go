package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CrawlJobHandler struct {
	crawlJobService domain.CrawlJobService
}

func NewHandler(crawlJobService domain.CrawlJobService) *CrawlJobHandler {
	return &CrawlJobHandler{
		crawlJobService: crawlJobService,
	}
}

func (h *CrawlJobHandler) Create(c *gin.Context) {
	var req dto.CreateCrawlJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	job, err := h.crawlJobService.Create(req.URL)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(201, dto.NewCreateCrawlJobResponse(job))
}
