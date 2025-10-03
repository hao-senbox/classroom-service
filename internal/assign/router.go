package assign

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *AssignHandler) {
	assginGroup := r.Group("/api/v1/admin/classrooms", middleware.Secured())
	{
		assginGroup.POST("/assigns", handler.AssignSlot)
		assginGroup.POST("/remove/assigns", handler.UnAssignSlot)

		
		// Assignment Template
		assginGroup.POST("/assignment-templates", handler.CreateAssignmentTemplate)
	}
}
