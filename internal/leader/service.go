package leader

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaderService interface {
	AddLeader(c *gin.Context, req *CreateLeaderRequest) error
}

type leaderService struct {
	LeaderRepository LeaderRepository
}

func NewLeaderService(leaderRepository LeaderRepository) LeaderService {
	return &leaderService{
		LeaderRepository: leaderRepository,
	}
}

func (s *leaderService) AddLeader(c *gin.Context, req *CreateLeaderRequest) error {

	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	objClassroomID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	data := &Leader{
		ID:          primitive.NewObjectID(),
		UserID:      req.UserID,
		ClassRoomID: objClassroomID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.LeaderRepository.CreateLeader(c, data)
}
