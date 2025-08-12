package class

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassRepository interface {
	CreateManyAssignments(ctx context.Context, assignments []*TeacherStudentAssignment) error
	GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error)
	GetAssgin(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error)
	FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error)
	UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error
	DeleteAssgin(ctx context.Context, classroomID primitive.ObjectID, index int) error
	CreateSystemNotification(ctx context.Context, system *SystemConfig) error
	GetFirstSystemNotification(ctx context.Context) (*SystemConfig, error)
	UpdateSystemNotification(ctx context.Context, system *SystemConfig) error
	FindAssignNotTeacher(ctx context.Context) ([]*TeacherStudentAssignment, error)
	FindAssign(ctx context.Context, classroomID primitive.ObjectID, index int, date *time.Time) (*TeacherStudentAssignment, error)
	GetAssignmentsByClassID(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) (bool, error)
	GetAssignmentsByClass(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) ([]*TeacherStudentAssignment, error)
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context) ([]*Notification, error)
	ReadNotification(ctx context.Context, id primitive.ObjectID) error
	CreateLeader(ctx context.Context, leader *Leader) error
	GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID) (*Leader, error)
}

type classRepository struct {
	assginCollection       *mongo.Collection
	systemConfigCollection *mongo.Collection
	notificationCollection *mongo.Collection
	leaderCollection       *mongo.Collection
}

func NewClassRepository(assginCollection, systemConfigCollection, notificationCollection, leaderCollection *mongo.Collection) ClassRepository {
	return &classRepository{
		assginCollection:       assginCollection,
		systemConfigCollection: systemConfigCollection,
		notificationCollection: notificationCollection,
		leaderCollection:       leaderCollection,
	}
}

func (r *classRepository) GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error) {

	var assgins []*TeacherStudentAssignment

	cursor, err := r.assginCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var assign TeacherStudentAssignment
		if err := cursor.Decode(&assign); err != nil {
			return nil, err
		}
		assgins = append(assgins, &assign)
	}

	return assgins, nil

}

func (r *classRepository) GetAssgin(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error) {

	var assign TeacherStudentAssignment

	if err := r.assginCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&assign); err != nil {
		return nil, err
	}

	return &assign, nil

}

func (r *classRepository) FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error) {

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

func (r *classRepository) UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error {
	_, err := r.assginCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": assign})
	return err
}

func (r *classRepository) DeleteAssgin(ctx context.Context, classroomID primitive.ObjectID, index int) error {

	filter := bson.M{
		"class_room_id": classroomID,
		"index":         index,
	}

	update := bson.M{
		"$set": bson.M{
			"student_id":      nil,
			"teacher_id":      nil,
			"is_notification": false,
		},
	}

	_, err := r.assginCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *classRepository) CreateSystemNotification(ctx context.Context, system *SystemConfig) error {
	_, err := r.systemConfigCollection.InsertOne(ctx, system)
	return err
}

func (r *classRepository) GetFirstSystemNotification(ctx context.Context) (*SystemConfig, error) {

	var system SystemConfig

	err := r.systemConfigCollection.FindOne(ctx, bson.M{}).Decode(&system)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &system, nil
}

func (r *classRepository) UpdateSystemNotification(ctx context.Context, system *SystemConfig) error {
	_, err := r.systemConfigCollection.UpdateOne(ctx, bson.M{}, bson.M{"$set": system})
	return err
}

func (r *classRepository) FindAssignNotTeacher(ctx context.Context) ([]*TeacherStudentAssignment, error) {

	var assgins []*TeacherStudentAssignment

	filter := bson.M{
		"teacher_id":      nil,
		"is_notification": false,
		"student_id":      bson.M{"$ne": nil},
	}

	cursor, err := r.assginCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var assign TeacherStudentAssignment
		if err := cursor.Decode(&assign); err != nil {
			return nil, err
		}
		assgins = append(assgins, &assign)
	}

	return assgins, nil

}

func (r *classRepository) CreateNotification(ctx context.Context, notification *Notification) error {
	_, err := r.notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

func (r *classRepository) GetNotifications(ctx context.Context) ([]*Notification, error) {

	var notifications []*Notification

	filter := bson.M{}

	cursor, err := r.notificationCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notification Notification
		if err := cursor.Decode(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (r *classRepository) ReadNotification(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.notificationCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_processed": true}})
	return err
}

func (r *classRepository) GetAssignmentsByClassID(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) (bool, error) {

	filter := bson.M{"class_room_id": classroomID}

	if date != nil {
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := start.Add(24 * time.Hour)

		filter["created_at"] = bson.M{
			"$gte": start,
			"$lt":  end,
		}
	}

	count, err := r.assginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *classRepository) CreateManyAssignments(ctx context.Context, assigns []*TeacherStudentAssignment) error {
	docs := make([]interface{}, len(assigns))
	for i, a := range assigns {
		docs[i] = a
	}
	_, err := r.assginCollection.InsertMany(ctx, docs)
	return err
}

func (r *classRepository) GetAssignmentsByClass(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) ([]*TeacherStudentAssignment, error) {

	filter := bson.M{"class_room_id": classroomID}

	if date != nil {
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := start.Add(24 * time.Hour)

		filter["created_at"] = bson.M{
			"$gte": start,
			"$lt":  end,
		}
	}

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

func (r *classRepository) FindAssign(ctx context.Context, classroomID primitive.ObjectID, index int, date *time.Time) (*TeacherStudentAssignment, error) {

	filter := bson.M{
		"class_room_id": classroomID,
		"index":         index,
	}

	if date != nil {
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := start.Add(24 * time.Hour)

		filter["created_at"] = bson.M{
			"$gte": start,
			"$lt":  end,
		}
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

func (r *classRepository) CreateLeader(ctx context.Context, leader *Leader) error {

	filter := bson.M{
		"class_room_id": leader.ClassRoomID,
	}

	_, err := r.leaderCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	_, err = r.leaderCollection.InsertOne(ctx, leader)
	if err != nil {
		return err
	}

	return nil

}

func (r *classRepository) GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID) (*Leader, error) {

	filter := bson.M{
		"class_room_id": classroomID,
	}

	var leader Leader
	err := r.leaderCollection.FindOne(ctx, filter).Decode(&leader)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &leader, nil
}
