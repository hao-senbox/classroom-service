package leader

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaderService interface {
	AddLeader(c *gin.Context, req *CreateLeaderRequest) error
	DeleteLeader(c *gin.Context, req *DeleteLeaderRequest) error
	// Leader Template
	CreateLeaderTemplate(c *gin.Context, req *CreateLeaderRequest) error
	DeleteLeaderTemplate(c *gin.Context, req *DeleteLeaderRequest) error
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

	if req.Owner == (Owner{}) {
		return fmt.Errorf("owner is required")
	}

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	objClassroomID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	data := &Leader{
		ID:          primitive.NewObjectID(),
		Owner:       req.Owner,
		ClassRoomID: objClassroomID,
		Date:        dateParse,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.LeaderRepository.CreateLeader(c, data)
}

func (s *leaderService) DeleteLeader(c *gin.Context, req *DeleteLeaderRequest) error {

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	objClassroomID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	return s.LeaderRepository.DeleteLeader(c, objClassroomID, &dateParse)
}

func (s *leaderService) CreateLeaderTemplate(c *gin.Context, req *CreateLeaderRequest) error {

	if req.Owner == (Owner{}) {
		return fmt.Errorf("owner is required")
	}

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	objClassroomID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	data := &LeaderTemplate{
		ID:          primitive.NewObjectID(),
		Owner:       &req.Owner,
		ClassRoomID: objClassroomID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.LeaderRepository.CreateLeaderTemplate(c, data)

}

func (s *leaderService) DeleteLeaderTemplate(c *gin.Context, req *DeleteLeaderRequest) error {

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	objClassroomID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	return s.LeaderRepository.DeleteLeaderTemplate(c, objClassroomID)

}