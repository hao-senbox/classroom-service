package leader

import (
	"classroom-service/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaderHandler struct {
	LeaderService LeaderService
}

func NewLeaderHandler(leaderService LeaderService) *LeaderHandler {
	return &LeaderHandler{
		LeaderService: leaderService,
	}
}

func (r *LeaderHandler) AddLeader(c *gin.Context) {

	var req CreateLeaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return 
	}

	err := r.LeaderService.AddLeader(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Add Leader Successfully", nil)

}

func (r *LeaderHandler) DeleteLeader(c *gin.Context) {

	var req DeleteLeaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return 
	}

	err := r.LeaderService.DeleteLeader(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Delete Leader Successfully", nil)
	
}

func (r *LeaderHandler) CreateLeaderTemplate(c *gin.Context) {

	var req CreateLeaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return 
	}

	err := r.LeaderService.CreateLeaderTemplate(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Create Leader Template Successfully", nil)

}

func (r *LeaderHandler) DeleteLeaderTemplate(c *gin.Context) {

	var req DeleteLeaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return 
	}

	err := r.LeaderService.DeleteLeaderTemplate(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, "INVALID_REQUEST")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Delete Leader Template Successfully", nil)
	
}