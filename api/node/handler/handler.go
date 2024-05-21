package handler

import (
	"crawlquery/api/domain"
	"crawlquery/api/dto"
	"crawlquery/api/errorutil"
	"net/http"
	"strconv"

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

	account := c.MustGet("account").(*domain.Account)
	node, err := h.nodeService.Create(account.ID, req.Hostname, req.Port)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(201, dto.NewCreateNodeResponse(node))
}

func (h *NodeHandler) ListByShardID(c *gin.Context) {

	shardID := c.Param("shardID")
	shardIdUint, err := strconv.ParseUint(shardID, 10, 64)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	nodes, err := h.nodeService.ListByShardID(uint(shardIdUint))

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(200, dto.NewListNodesByShardResponse(nodes))
}

func (h *NodeHandler) ListByAccountID(c *gin.Context) {
	account := c.MustGet("account").(*domain.Account)
	nodes, err := h.nodeService.ListByAccountID(account.ID)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(200, dto.NewListNodesResponse(nodes))
}

func (h *NodeHandler) Auth(c *gin.Context) {
	var req dto.AuthenticateNodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorutil.HandleGinError(c, err, http.StatusBadRequest)
		return
	}

	node, err := h.nodeService.Auth(req.Key)

	if err != nil {
		errorutil.HandleGinError(c, err, http.StatusUnauthorized)
		return
	}

	c.JSON(200, dto.NewAuthenticateNodeResponse(node))
}
