package leader

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *LeaderHandler) {
	leaderGroup := r.Group("/api/v1/admin/classrooms", middleware.Secured())
	{
		leaderGroup.POST("/leader", handler.AddLeader)
		leaderGroup.POST("/remove/leader", handler.DeleteLeader)

		// Leader Template
		leaderGroup.POST("/leader-templates", handler.CreateLeaderTemplate)
		leaderGroup.POST("/remove/leader-templates", handler.DeleteLeaderTemplate)
	}
}
