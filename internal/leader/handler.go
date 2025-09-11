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
