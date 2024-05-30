package handler

import (
	"crawlquery/html/dto"
	"crawlquery/pkg/util"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storagePath string
}

func NewHandler(storagePath string) *Handler {
	return &Handler{storagePath: storagePath}
}

func (h Handler) GetPage(c *gin.Context) {

	pageID := c.Param("pageID")

	if !util.ValidatePageID(pageID) {
		c.JSON(404, dto.ErrorResponse{Error: "invalid page ID"})
		return
	}

	// try to read page
	data, err := os.ReadFile(h.storagePath + "/" + pageID)

	if err != nil {
		c.JSON(404, dto.ErrorResponse{Error: "page not found"})
		return
	}

	c.JSON(200, dto.GetPageResponse{HTML: data})
}

func (h Handler) StorePage(c *gin.Context) {

	var req dto.StorePageRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	if !util.ValidateHash(req.Hash) {
		c.JSON(400, dto.ErrorResponse{Error: "invalid page ID"})
		return
	}

	if err := os.WriteFile(h.storagePath+"/"+req.Hash, req.HTML, 0644); err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "failed to store page"})
		return
	}

	c.JSON(201, dto.StorePageResponse{
		Success: true,
	})
}
