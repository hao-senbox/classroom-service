package room

import (
	"classroom-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *RoomHandler) {
	room := r.Group("api/v1/room").Use(middleware.Secured())
	{
		room.POST("", handler.CreateRoom)
		room.GET("", handler.GetRooms)
		room.GET("/:id", handler.GetRoom)
		room.PUT("/:id", handler.UpdateRoom)
		room.DELETE("/:id", handler.DeleteRoom)
	}
	assgin := r.Group("api/v1/assgin").Use(middleware.Secured())
	{
		assgin.POST("", handler.CreateAssgin)
		assgin.GET("", handler.GetAssgins)
		assgin.GET("/:id", handler.GetAssgin)
		assgin.PUT("/:id", handler.UpdateAssgin)
		assgin.DELETE("/:id", handler.DeleteAssgin)
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