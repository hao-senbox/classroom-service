package class

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	ClassService ClassService
}

func NewClassHandler(classService ClassService) *ClassHandler {
	return &ClassHandler{
		ClassService: classService,
	}
}

func (h *ClassHandler) CreateClass(c *gin.Context) {

	var req CreateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	id, err := h.ClassService.CreateClass(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", id)
}

func (h *ClassHandler) GetClasses(c *gin.Context) {

	date := c.Query("date")

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	rooms, err := h.ClassService.GetClasses(ctx, date)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", rooms)
}

func (h *ClassHandler) AddLeader(c *gin.Context) {

	var req AddLeaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.ClassService.AddLeader(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}

func (h *ClassHandler) GetAssgins(c *gin.Context) {

	assgins, err := h.ClassService.GetAssgins(c)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", assgins)

}

func (h *ClassHandler) GetAssgin(c *gin.Context) {

	index := c.Param("index")
	if index == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("index is required"), "INVALID_REQUEST")
		return
	}

	classroomID := c.Query("classroom_id")
	if classroomID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("classroom_id is required"), "INVALID_REQUEST")
		return
	}

	date := c.Query("date")
	if date == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("date is required"), "INVALID_REQUEST")
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	assgin, err := h.ClassService.GetAssgin(ctx, classroomID, index, date)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", assgin)

}

func (h *ClassHandler) CreateAssgin(c *gin.Context) {

	var req UpdateAssginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.ClassService.CreateAssgin(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}

func (h *ClassHandler) DeleteAssgin(c *gin.Context) {

	var req UpdateAssginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.ClassService.DeleteAssgin(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *ClassHandler) CreateSystemNotification(c *gin.Context) {

	var req CreateSystemNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	err := h.ClassService.CreateSystemNotification(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *ClassHandler) GetSystemNotification(c *gin.Context) {

	system, err := h.ClassService.GetSystemNotification(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", system)
}

func (h *ClassHandler) UpdateSystemNotification(c *gin.Context) {

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

	err := h.ClassService.UpdateSystemNotification(c, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}

func (h *ClassHandler) GetNotifications(c *gin.Context) {

	notifications, err := h.ClassService.GetNotifications(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", notifications)

}

func (h *ClassHandler) ReadNotification(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	err := h.ClassService.ReadNotification(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}
