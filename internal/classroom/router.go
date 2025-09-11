package classroom

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ClassroomHandler) {
	classroomGroup := r.Group("/api/v1/classrooms", middleware.Secured())
	{
		classroomGroup.POST("", handler.CreateClassroom)
		// classroomGroup.GET("", handler.GetClassrooms)
		// classroomGroup.GET("/:id", handler.GetClassroom)
		// classroomGroup.PUT("/:id", handler.UpdateClassroom)
		// classroomGroup.DELETE("/:id", handler.DeleteClassroom)
	}
	apiGatewayClassroomGroup := r.Group("/api/v1/gateway", middleware.Secured())
	{
		apiGatewayClassroomGroup.GET("/classroom", handler.GetClassroomsByUserID)
	}
}
