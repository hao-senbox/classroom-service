package leader

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *LeaderHandler) {
	leaderGroup := r.Group("/api/v1/leader", middleware.Secured())
	{
		leaderGroup.POST("", handler.AddLeader)
	}
}
