package storage

import (
	"context"
	"fmt"
	"log"

	pb "budget-service/genproto"
	"budget-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationService struct {
	db *mongo.Database
}

func NewNotificationService(db *mongo.Database) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) CreateNotification(req model.Send) error {
	coll := s.db.Collection("notifications")
	// Generate a new ObjectID for the notification
	objID := primitive.NewObjectID()

	_, err := coll.InsertOne(context.Background(), bson.M{
		"_id":     objID, // Use ObjectID for _id
		"user_id":  req.UserId,
		"message": req.Message,
	})
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return err
	}
	return nil
}

// GetNotification retrieves a notification by user_id
func (s *NotificationService) GetNotification(req *pb.GetNotificationByidRequest) (*pb.GetNotificationByidResponse, error) {
	coll := s.db.Collection("notifications")
	var notificationData struct {
		ID      primitive.ObjectID `bson:"_id"` // BSON tag for ID field
		UserId  string             `bson:"user_Id"`
		Message string             `bson:"message"`
	}
	err := coll.FindOne(context.Background(), bson.M{"user_id": req.UserId}).Decode(&notificationData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("notification not found")
		}
		log.Printf("Failed to retrieve notification: %v", err)
		return nil, err
	}

	return &pb.GetNotificationByidResponse{
		UserId:  notificationData.UserId,
		Message: notificationData.Message,
	}, nil
}

// DeleteNotification deletes a notification by user_id
func (s *NotificationService) DeleteNotification(req *pb.GetNotificationByidRequest) (*pb.NotificationsResponse, error) {
	coll := s.db.Collection("notifications")
	res, err := coll.DeleteOne(context.Background(), bson.M{"user_Id": req.UserId})
	if err != nil {
		log.Printf("Failed to delete notification: %v", err)
		return &pb.NotificationsResponse{
			Message: "Error while deleting notification",
			Success: false,
		}, err
	}

	if res.DeletedCount == 0 {
		return &pb.NotificationsResponse{
			Message: "No notification found to delete",
			Success: false,
		}, nil
	}

	return &pb.NotificationsResponse{
		Message: "Notification deleted successfully",
		Success: true,
	}, nil
}

// ListNotification lists all notifications
func (s *NotificationService) ListNotification(req *pb.Void) (*pb.ListNotificationResponse, error) {
	coll := s.db.Collection("notifications")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notifications []*pb.GetNotificationByidResponse
	for cursor.Next(context.Background()) {
		var notificationData struct {
			ID      primitive.ObjectID `bson:"_id"`
			UserId  string             `bson:"user_id"`
			Message string             `bson:"message"`
		}
		if err := cursor.Decode(&notificationData); err != nil {
			return nil, err
		}
		notifications = append(notifications, &pb.GetNotificationByidResponse{
			UserId:  notificationData.UserId,
			Message: notificationData.Message,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListNotificationResponse{Notifications: notifications}, nil
}
