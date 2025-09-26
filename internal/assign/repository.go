package assign

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssignRepository interface {
	CreateAssignment(ctx context.Context, assign *TeacherStudentAssignment) error
	CheckDuplicateAssignmentForDate(ctx context.Context, classroomID primitive.ObjectID, date time.Time, studentID, teacherID string) (bool, error)
	GetAssignmentBySlotAndDate(ctx context.Context, classroomID primitive.ObjectID, slotNumber int, date *time.Time) (*TeacherStudentAssignment, error)
	UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error
	GetAssignmentsByClassroomAndDate(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) ([]*TeacherStudentAssignment, error)
	CountAssignedSlotsTotal(ctx context.Context, classroomID primitive.ObjectID) (int, error)
	GetAssignmentsByClassroomID(ctx context.Context, classroomID primitive.ObjectID) ([]*TeacherStudentAssignment, error)
}

type assignRepository struct {
	assginCollection *mongo.Collection
}

func NewAssignRepository(assginCollection *mongo.Collection) AssignRepository {
	return &assignRepository{
		assginCollection: assginCollection,
	}
}

func (r *assignRepository) CreateAssignment(ctx context.Context, assign *TeacherStudentAssignment) error {

	_, err := r.assginCollection.InsertOne(ctx, assign)
	if err != nil {
		return err
	}

	return err
}

func (r *assignRepository) GetAssignmentBySlotAndDate(ctx context.Context, classroomID primitive.ObjectID, slotNumber int, date *time.Time) (*TeacherStudentAssignment, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	filter := bson.M{
		"class_room_id": classroomID,
		"slot_number":   slotNumber,
		"assign_date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	var assign TeacherStudentAssignment
	err := r.assginCollection.FindOne(ctx, filter).Decode(&assign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &assign, nil

}

func (r *assignRepository) CheckDuplicateAssignmentForDate(ctx context.Context, classroomID primitive.ObjectID, date time.Time, studentID, teacherID string) (bool, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	filter := bson.M{
		"class_room_id": classroomID,
		"student_id":    studentID,
		"teacher_id":    teacherID,
		"assign_date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	count, err := r.assginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *assignRepository) UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error {
	_, err := r.assginCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": assign})
	return err
}

func (r *assignRepository) GetAssignmentsByClassroomAndDate(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) ([]*TeacherStudentAssignment, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	filter := bson.M{
		"class_room_id": classroomID,
		"assign_date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	fmt.Printf("filter: %v\n", filter)

	cursor, err := r.assginCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*TeacherStudentAssignment
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}

func (r *assignRepository) CountAssignedSlotsTotal(ctx context.Context, classroomID primitive.ObjectID) (int, error) {
	count, err := r.assginCollection.CountDocuments(ctx, bson.M{"class_room_id": classroomID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *assignRepository) GetAssignmentsByClassroomID(ctx context.Context, classroomID primitive.ObjectID) ([]*TeacherStudentAssignment, error) {

	filter := bson.M{"class_room_id": classroomID}

	cursor, err := r.assginCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []*TeacherStudentAssignment
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}
