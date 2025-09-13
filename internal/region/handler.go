package region

import (
	"classroom-service/helper"
	"classroom-service/pkg/constants"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	RegionService RegionService
}

func NewRegionHandler(regionService RegionService) *RegionHandler {
	return &RegionHandler{
		RegionService: regionService,
	}
}

func (h *RegionHandler) CreateRegion(c *gin.Context) {

	var req CreateRegionRequest

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

	id, err := h.RegionService.CreateRegion(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Created Region Successfully", id)

}

func (h *RegionHandler) GetRegions(c *gin.Context) {

	organizationID := c.Query("organization_id")
	if organizationID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("organization_id is required"), "INVALID_REQUEST")
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

	regions, err := h.RegionService.GetAllRegions(ctx, organizationID, date)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get Regions Successfully", regions)

}

func (h *RegionHandler) GetRegion(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
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

	region, err := h.RegionService.GetRegion(ctx, id, date)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get Region Successfully", region)

}

func (h *RegionHandler) UpdateRegion(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	var req UpdateRegionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	err := h.RegionService.UpdateRegion(ctx, id, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Update Region Successfully", nil)

}

func (h *RegionHandler) DeleteRegion(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), "INVALID_REQUEST")
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), "UNAUTHORIZED")
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	err := h.RegionService.DeleteRegion(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Delete Region Successfully", nil)

}
