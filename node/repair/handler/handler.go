package handler

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repairService domain.RepairService
}

func NewHandler(repairService domain.RepairService) *Handler {
	return &Handler{
		repairService: repairService,
	}
}

func (h *Handler) GetIndexMetas(c *gin.Context) {

	var req dto.GetIndexMetasRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	metas, err := h.repairService.GetIndexMetas(req.PageIDs)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var dtoIndexMetas []dto.IndexMeta

	for _, meta := range metas {
		dtoIndexMetas = append(dtoIndexMetas, dto.IndexMeta{
			PageID:        string(meta.PageID),
			PeerID:        string(meta.PeerID),
			LastIndexedAt: meta.LastIndexedAt,
		})
	}

	c.JSON(200, dto.GetIndexMetasResponse{
		IndexMetas: dtoIndexMetas,
	})
}
