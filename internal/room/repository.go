package room

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *ClassRoom) (string, error)
	GetRooms(ctx context.Context) ([]*ClassRoom, error)
	GetRoom(ctx context.Context, id primitive.ObjectID) (*ClassRoom, error)
	UpdateRoom(ctx context.Context, id primitive.ObjectID, room *ClassRoom) error
	DeleteRoom(ctx context.Context, id primitive.ObjectID) error
	CreateAssgin(ctx context.Context, assign *TeacherStudentAssignment) (string, error) 
	GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error)
	GetAssgin(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error)
	FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error)
	UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error
	DeleteAssgin(ctx context.Context, id primitive.ObjectID) error
	CreateSystemNotification(ctx context.Context, system *SystemConfig) error
	GetFirstSystemNotification(ctx context.Context) (*SystemConfig, error)
	UpdateSystemNotification(ctx context.Context, system *SystemConfig) error
	FindAssignNotTeacher(ctx context.Context) ([]*TeacherStudentAssignment, error)
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context) ([]*Notification, error)
	ReadNotification(ctx context.Context, id primitive.ObjectID) error
}

type roomRepository struct {
	roomCollection *mongo.Collection
	assginCollection *mongo.Collection
	systemConfigCollection *mongo.Collection
	notificationCollection *mongo.Collection
}

func NewRoomRepository(roomCollection, assginCollection, systemConfigCollection, notificationCollection *mongo.Collection) RoomRepository {
	return &roomRepository{
		roomCollection: roomCollection,
		assginCollection: assginCollection,
		systemConfigCollection: systemConfigCollection,
		notificationCollection: notificationCollection,
	}
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *ClassRoom) (string, error) {
	
	result, err := r.roomCollection.InsertOne(ctx, room)
	if err != nil {
		return "", err
	}
	
	oid := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

func (r *roomRepository) GetRooms(ctx context.Context) ([]*ClassRoom, error) {

	var rooms []*ClassRoom

	cursor, err := r.roomCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var room ClassRoom
		if err := cursor.Decode(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, nil
	
}

func (r *roomRepository) GetRoom(ctx context.Context, id primitive.ObjectID) (*ClassRoom, error) {

	var room ClassRoom

	if err := r.roomCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) UpdateRoom(ctx context.Context, id primitive.ObjectID, room *ClassRoom) error {
	_, err := r.roomCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": room})
	return err
}

func (r *roomRepository) DeleteRoom(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.roomCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *roomRepository) CreateAssgin(ctx context.Context, assign *TeacherStudentAssignment) (string, error) {

	result, err := r.assginCollection.InsertOne(ctx, assign)
	if err != nil {
		return "", err
	}
	
	oib := result.InsertedID.(primitive.ObjectID)

	return oib.Hex(), nil

}	

func (r *roomRepository) GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error) {

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

func (r *roomRepository) GetAssgin(ctx context.Context, id primitive.ObjectID) (*TeacherStudentAssignment, error) {

	var assign TeacherStudentAssignment

	if err := r.assginCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&assign); err != nil {
		return nil, err
	}

	return &assign, nil

}

func (r *roomRepository) FindDuplicate(ctx context.Context, classroomID primitive.ObjectID, studentID, teacherID string) (bool, error) {

	filter := bson.M{
		"class_room_id": classroomID,
		"student_id": studentID,
		"teacher_id": teacherID,
	}

	count, err := r.assginCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
	
}

func (r *roomRepository) UpdateAssgin(ctx context.Context, id primitive.ObjectID, assign *TeacherStudentAssignment) error {
	_, err := r.assginCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": assign})
	return err
}

func (r *roomRepository) DeleteAssgin(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.assginCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *roomRepository) CreateSystemNotification(ctx context.Context, system *SystemConfig) error {
	_, err := r.systemConfigCollection.InsertOne(ctx, system)
	return err
}

func (r *roomRepository) GetFirstSystemNotification(ctx context.Context) (*SystemConfig, error) {
	
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

func (r *roomRepository) UpdateSystemNotification(ctx context.Context, system *SystemConfig) error {
	_, err := r.systemConfigCollection.UpdateOne(ctx, bson.M{}, bson.M{"$set": system})
	return err
}

func (r *roomRepository) FindAssignNotTeacher(ctx context.Context) ([]*TeacherStudentAssignment, error) {

	var assgins []*TeacherStudentAssignment

	filter := bson.M{"teacher_id": nil, "is_notification": false}

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

func (r *roomRepository) CreateNotification(ctx context.Context, notification *Notification) error {
	_, err := r.notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) GetNotifications(ctx context.Context) ([]*Notification, error) {

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

func (r *roomRepository) ReadNotification(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.notificationCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_processed": true}})
	return err
}