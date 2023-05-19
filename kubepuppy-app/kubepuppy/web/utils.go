package openapi

import (
	"github.com/gin-gonic/gin"
	"kubepuppy/app"
)

const ContextClusterKey = "alanpinder.com/kubepuppy/cluster"

func GetClusterFromRequest(c *gin.Context) *app.Cluster {
	clusterAsAny, ok := c.Get(ContextClusterKey)

	if !ok {
		panic("cluster not supplied")
	}

	return clusterAsAny.(*app.Cluster)
}

func SetClusterIntoRequest(c *gin.Context, cluster *app.Cluster) {
	c.Set(ContextClusterKey, cluster)
}
