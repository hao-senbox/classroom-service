package assign

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *AssignHandler) {
	assginGroup := r.Group("/api/v1/assgin", middleware.Secured())
	{
		// assginGroup.GET("", handler.GetAssgins)
		assginGroup.POST("", handler.AssignSlot)
		// assginGroup.DELETE("", handler.DeleteAssgin)
		// assginGroup.GET("/:index", handler.GetAssgin)
	}
}
