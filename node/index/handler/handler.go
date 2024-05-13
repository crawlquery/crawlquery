package handler

import (
	"crawlquery/node/index"

	"github.com/gin-gonic/gin"
)

type IndexHandler struct {
	idx *index.Index
}

func NewHandler(idx *index.Index) *IndexHandler {
	return &IndexHandler{
		idx: idx,
	}
}

func (sh *IndexHandler) Search(c *gin.Context) {

	res, err := sh.idx.Search(c.Query("q"))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"results": res,
	})
}
