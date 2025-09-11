package assign

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssignHandler struct {
	AssignService AssignService
}

func NewAssignHandler(assignService AssignService) *AssignHandler {
	return &AssignHandler{
		AssignService: assignService,
	}
}

func (h *AssignHandler) AssignSlot(c *gin.Context) {

	var req UpdateAssginRequest

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

	err := h.AssignService.AssignSlot(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)
}
