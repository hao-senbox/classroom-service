package assign

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssignRepository interface {
	CreateManyAssignments(ctx context.Context, assigns []*TeacherStudentAssignment) error
	FindAssign(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error)
	FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error)
	UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error
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

func (r *assignRepository) CreateManyAssignments(ctx context.Context, assigns []*TeacherStudentAssignment) error {

	docs := make([]interface{}, len(assigns))
	for i, a := range assigns {
		docs[i] = a
	}

	_, err := r.assginCollection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil

}

func (r *assignRepository) FindAssign(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error) {

	var assign TeacherStudentAssignment
	err := r.assginCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&assign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &assign, nil
}

func (r *assignRepository) FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error) {

	filter := bson.M{
		"class_room_id": classroomID,
		"student_id":    studentID,
		"teacher_id":    teacherID,
	}

	count, err := r.assginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *assignRepository) UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error {

	_, err := r.assginCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": assign})
	if err != nil {
		return err
	}

	return nil

}

func (r *assignRepository) GetAssignmentsByClassroomID(ctx context.Context, classroomID primitive.ObjectID) ([]*TeacherStudentAssignment, error) {

	cursor, err := r.assginCollection.Find(ctx, bson.M{"class_room_id": classroomID})
	if err != nil {
		return nil, err
	}

	var results []*TeacherStudentAssignment
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}