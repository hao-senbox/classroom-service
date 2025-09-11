package class

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ClassHandler) {
	class := r.Group("api/v1/room").Use(middleware.Secured())
	{
		class.POST("", handler.CreateClass)
		class.PUT("/add/leader", handler.AddLeader)
	}
	assgin := r.Group("api/v1/assgin").Use(middleware.Secured())
	{
		assgin.GET("", handler.GetAssgins)
		assgin.PUT("", handler.CreateAssgin)
		assgin.DELETE("", handler.DeleteAssgin)
		assgin.GET("/:index", handler.GetAssgin)
	}
	notification := r.Group("api/v1/notification").Use(middleware.Secured())
	{
		notification.GET("", handler.GetNotifications)
		notification.PUT("/:id", handler.ReadNotification)
		notification.POST("delay", handler.CreateSystemNotification)
		notification.GET("delay", handler.GetSystemNotification)
		notification.PUT("delay/:id", handler.UpdateSystemNotification)
	}
}
