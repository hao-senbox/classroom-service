package assign

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *AssignHandler) {
	assginGroup := r.Group("/api/v1/admin/classrooms", middleware.Secured())
	{
		// assginGroup.GET("", handler.GetAssgins)
		assginGroup.POST("/assigns", handler.AssignSlot)
		assginGroup.POST("/remove/assigns", handler.UnAssignSlot)
		// assginGroup.GET("/:index", handler.GetAssgin)
	}
}
