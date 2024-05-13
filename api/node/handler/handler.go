package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService domain.NodeService
}

func NewHandler(nodeService domain.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
	}
}

func (h *NodeHandler) Create(c *gin.Context) {
	var req dto.CreateNodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	accountID := c.MustGet("account_id").(string)
	node, err := h.nodeService.Create(accountID, req.Hostname, req.Port)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(201, dto.NewCreateNodeResponse(node))
}
