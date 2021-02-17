package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var groups map[string]*Group = nil

func getForRemotePeer(c *gin.Context) {
	groupName := c.Param("group")
	key := c.Param("key")
	if group, ok := groups[groupName]; ok {
		if value := group.Get(key); value != nil {
			c.String(http.StatusOK, string(value))
			return
		}
	}
	c.String(http.StatusNotFound, "")
}

func StartServer(address string, _groups map[string]*Group) {
	groups = _groups
	r := gin.Default()
	r.GET("/cache/:group/:key", getForRemotePeer)
	r.Run(address)
}
