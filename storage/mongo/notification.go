package storage

import (
	ctx "context"
	"fmt"
	"log"

	pb "budget-service/genproto"
	"budget-service/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
	id := uuid.NewString()
	_, err := coll.InsertOne(ctx.Background(), bson.M{
		"id":      id,
		"UserId":  req.UserId,
		"message": req.Message,
	})
	if err != nil {
		log.Printf("Failed to create account: %v", err)
		return err
	}
	return nil
}

// GetAccountByid retrieves a notification by user_id
func (s *NotificationService) GetNotification(req *pb.GetNotificationByidRequest) (*pb.GetNotificationByidResponse, error) {
	coll := s.db.Collection("notifications")
	var notification model.Send
	err := coll.FindOne(ctx.Background(), bson.M{"user_id": req.UserId}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("notification not found")
		}
		log.Printf("Failed to retrieve notification: %v", err)
		return nil, err
	}

	return &pb.GetNotificationByidResponse{
		UserId:  notification.UserId,
		Message: notification.Message,
	}, nil
}

// DeleteAccount deletes a notification by user_id
func (s *NotificationService) DeleteNotification(req *pb.GetNotificationByidRequest) (*pb.NotificationsResponse, error) {
	coll := s.db.Collection("notifications")
	res, err := coll.DeleteOne(ctx.Background(), bson.M{"user_id": req.UserId})
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

// ListAccounts lists all notifications
func (s *NotificationService) ListNotification(req *pb.Void) (*pb.ListNotificationResponse, error) {
	coll := s.db.Collection("notifications")
	cursor, err := coll.Find(ctx.Background(), bson.M{})
	if err != nil {
		log.Printf("Failed to list notifications: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx.Background())

	var notifications []*pb.GetNotificationByidResponse
	for cursor.Next(ctx.Background()) {
		var notification model.Send
		if err := cursor.Decode(&notification); err != nil {
			log.Printf("Failed to decode notification: %v", err)
			return nil, err
		}
		notifications = append(notifications, &pb.GetNotificationByidResponse{
			UserId:  notification.UserId,
			Message: notification.Message,
		})
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListNotificationResponse{Notifications: notifications}, nil
}
