package classroom

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassroomHandler struct {
	ClassroomService ClassroomService
}

func NewClassroomHandler(classroomService ClassroomService) *ClassroomHandler {
	return &ClassroomHandler{
		ClassroomService: classroomService,
	}
}

func (h *ClassroomHandler) CreateClassroom(c *gin.Context) {

	var req CreateClassroomRequest

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

	id, err := h.ClassroomService.CreateClassroom(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create Classroom Successfully", id)

}

func (h *ClassroomHandler) GetClassroomsByUserID(c *gin.Context) {

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

	classrooms, err := h.ClassroomService.GetClassroomsByUserID(ctx, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", classrooms)
}
