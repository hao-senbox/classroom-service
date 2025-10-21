package classroom

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ClassroomHandler) {
	classroomGroup := r.Group("/api/v1/admin/classrooms", middleware.Secured())
	{
		classroomGroup.POST("", handler.CreateClassroom)
		classroomGroup.GET("/:id", handler.GetClassroomByID)
		classroomGroup.PUT("/:id", handler.UpdateClassroom)

		// Classroom Template
		classroomGroup.GET("/template/:classroom_id", handler.GetClassroomByIDTemplate)
		classroomGroup.POST("/template", handler.CreateAssignmentByTemplate)

		// Classroom Assignment
		classroomGroup.GET("/teacher-assignments", handler.GetTeacherAssignments)

	}
	apiGatewayClassroomGroup := r.Group("/api/v1/gateway", middleware.Secured())
	{
		apiGatewayClassroomGroup.GET("/classrooms", handler.GetClassroomsByOrg)
		apiGatewayClassroomGroup.GET("/classrooms/teacher-assignments", handler.GetTeacherAssignmentsByClassroomID)
		apiGatewayClassroomGroup.GET(("/classrooms/term"), handler.GetStudentsByTermAndClassroomID)
		apiGatewayClassroomGroup.GET("/classrooms/template", handler.GetClassroomTemplateByClassroomID)
	}
}
