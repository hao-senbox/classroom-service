package classroom

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func (h *ClassroomHandler) GetClassroomByID(c *gin.Context) {

	start := c.Query("start")
	end := c.Query("end")

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	classroom, err := h.ClassroomService.GetClassroomByID(ctx, id, start, end, page, limit)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get Classroom Successfully", classroom)

}

func (h *ClassroomHandler) UpdateClassroom(c *gin.Context) {

	var req UpdateClassroomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ClassroomService.UpdateClassroom(ctx, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Update Classroom Successfully", nil)

}

func (h *ClassroomHandler) GetClassroomsByUserID(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	classrooms, err := h.ClassroomService.GetClassroomsByUserID(ctx)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get Classrooms Successfully", classrooms)
}

func (h *ClassroomHandler) GetClassroomByIDTemplate(c *gin.Context) {

	id := c.Param("classroom_id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	classroom, err := h.ClassroomService.GetClassroomByIDTemplate(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get Classroom Template Successfully", classroom)

}

func (h *ClassroomHandler) CreateAssignmentByTemplate(c *gin.Context) {

	var req CreateAssignmentByTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ClassroomService.CreateAssignmentByTemplate(ctx, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create Assignment Successfully", nil)

}

func (h *ClassroomHandler) GetTeacherAssignments(c *gin.Context) {

	organizationID := c.Query("organization_id")
	termID := c.Query("term_id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	assignments, err := h.ClassroomService.GetTeacherAssignments(ctx, userID.(string), organizationID, termID)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get teacher assignments successfully", assignments)

}

func (h *ClassroomHandler) GetTeacherAssignmentsByClassroomID(c *gin.Context) {

	classroomID := c.Query("classroom_id")
	if classroomID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("classroom_id is required"), "INVALID_REQUEST")
		return
	}

	teacherID := c.Query("teacher_id")
	if teacherID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("teacher_id is required"), "INVALID_REQUEST")
		return
	}

	termID := c.Query("term_id")
	if termID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("term_id is required"), "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	assignments, err := h.ClassroomService.GetTeacherAssignmentsByClassroomID(ctx, classroomID, teacherID, termID)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get teacher assignments successfully", assignments)

}

func (h *ClassroomHandler) GetStudentsByTermAndClassroomID(c *gin.Context) {

	termID := c.Query("term_id")
	if termID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("term_id is required"), "INVALID_REQUEST")
		return
	}

	classroomID := c.Query("classroom_id")
	if classroomID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("classroom_id is required"), "INVALID_REQUEST")
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	students, err := h.ClassroomService.GetStudentsByTermAndClassroomID(ctx, classroomID, termID)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get students successfully", students)
	
}