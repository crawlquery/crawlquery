package handler

import (
	"crawlquery/api/service"
	"crawlquery/pkg/domain"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(ns *service.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: ns,
	}
}

func (nh *NodeHandler) Add(c *gin.Context) {
	var req struct {
		Hostname string `json:"hostname" binding:"required"`
		Port     string `json:"port" binding:"required"`
		ShardID  int    `json:"shard_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := nh.nodeService.Create(&domain.Node{
		Hostname: req.Hostname,
		Port:     req.Port,
		ShardID:  domain.ShardID(req.ShardID),
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
