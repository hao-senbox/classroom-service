package room

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	RoomService RoomService
}

func NewRoomHandler(roomService RoomService) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	userID, exist := c.Get(constants.UserID)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	roomID, err := h.RoomService.CreateRoom(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", roomID)

}

func (h *RoomHandler) GetRooms(c *gin.Context) {

	rooms, err := h.RoomService.GetRooms(context.Background())

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", rooms)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	room, err := h.RoomService.GetRoom(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", room)

}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {

	var req UpdateRoomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.RoomService.UpdateRoom(c, &req, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.RoomService.DeleteRoom(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *RoomHandler) CreateAssgin(c *gin.Context) {

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	userID, exist := c.Get(constants.UserID)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	var req CreateAssginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	assignID, err := h.RoomService.CreateAssgin(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", assignID)

}

func (h *RoomHandler) GetAssgins(c *gin.Context) {

	assgins, err := h.RoomService.GetAssgins(c)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", assgins)

}

func (h *RoomHandler) GetAssgin(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	assgin, err := h.RoomService.GetAssgin(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", assgin)

}

func (h *RoomHandler) UpdateAssgin(c *gin.Context) {

	var req UpdateAssginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.RoomService.UpdateAssgin(c, &req, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}

func (h *RoomHandler) DeleteAssgin(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.RoomService.DeleteAssgin(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *RoomHandler) CreateSystemNotification(c *gin.Context) {

	var req CreateSystemNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.RoomService.CreateSystemNotification(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *RoomHandler) GetSystemNotification(c *gin.Context) {

	system, err := h.RoomService.GetSystemNotification(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", system)
}

func (h *RoomHandler) UpdateSystemNotification(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	var req UpdateSystemNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.RoomService.UpdateSystemNotification(c, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}

func (h *RoomHandler) GetNotifications(c *gin.Context) {

	notifications, err := h.RoomService.GetNotifications(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", notifications)

}

func (h *RoomHandler) ReadNotification(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.RoomService.ReadNotification(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}