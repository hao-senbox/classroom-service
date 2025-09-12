package region

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *RegionHandler) {
	regionGroup := r.Group("/api/v1/admin/classrooms/regions", middleware.Secured())
	{
		regionGroup.POST("", handler.CreateRegion)
		regionGroup.GET("", handler.GetRegions)
		regionGroup.GET("/:id", handler.GetRegion)
		regionGroup.PUT("/:id", handler.UpdateRegion)
		regionGroup.DELETE("/:id", handler.DeleteRegion)
	}
}
