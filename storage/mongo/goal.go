package storage

import (
	"context"
	"log"
	"time"

	pb "budget-service/genproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GoalStorage struct to handle goal operations in MongoDB
type GoalStorage struct {
	db *mongo.Database
}

// NewGoalStorage initializes a new GoalStorage
func NewGoalStorage(db *mongo.Database) *GoalStorage {
	return &GoalStorage{db: db}
}

// CreateGoal creates a new goal in the database
func (s *GoalStorage) CreateGoal(req *pb.CreateGoalRequest) (*pb.Responsee, error) {
	coll := s.db.Collection("goals")

	_, err := coll.InsertOne(context.Background(), bson.M{
		"id":             req.Id,
		"user_id":        req.UserId,
		"target_amount":  req.TargetAmount,
		"current_amount": req.CurrentAmount,
		"deadline":       req.Deadline,
		"status":         req.Status,
	})
	if err != nil {
		log.Printf("Failed to create goal: %v", err)
		return &pb.Responsee{Message: "Failed to create goal"}, err
	}

	return &pb.Responsee{Message: "Goal created successfully"}, nil
}

// ListGoals lists all goals
func (s *GoalStorage) ListGoals(req *pb.ListGoalsRequest) (*pb.ListGoalsResponse, error) {
	coll := s.db.Collection("goals")

	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.TargetAmount > 0 {
		filter["target_amount"] = req.TargetAmount
	}
	if req.CurrentAmount > 0 {
		filter["current_amount"] = req.CurrentAmount
	}
	if req.Deadline != "" {
		filter["deadline"] = req.Deadline
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to list goals: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var goalsList []*pb.GoalResponse
	for cursor.Next(context.Background()) {
		var goal pb.GoalResponse
		if err := cursor.Decode(&goal); err != nil {
			log.Printf("Failed to decode goal: %v", err)
			return nil, err
		}
		goalsList = append(goalsList, &goal)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListGoalsResponse{Goals: goalsList}, nil
}

// GetGoalById retrieves a goal by its ID
func (s *GoalStorage) GetGoalById(req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error) {
	coll := s.db.Collection("goals")

	var goal pb.GoalResponse
	err := coll.FindOne(context.Background(), bson.M{"id": req.GoalId}).Decode(&goal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No goal found with id: %v", req.GoalId)
			return nil, err
		}
		log.Printf("Failed to get goal by id: %v", err)
		return nil, err
	}

	return &goal, nil
}

// UpdateGoal updates a goal based on the provided request data
func (s *GoalStorage) UpdateGoal(req *pb.UpdateGoalRequest) (*pb.Responsee, error) {
	coll := s.db.Collection("goals")

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.TargetAmount > 0 {
		update["target_amount"] = req.TargetAmount
	}
	if req.CurrentAmount > 0 {
		update["current_amount"] = req.CurrentAmount
	}
	if req.Deadline != "" {
		update["deadline"] = req.Deadline
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	if len(update) == 0 {
		return &pb.Responsee{Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(context.Background(), bson.M{"id": req.GoalId}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update goal: %v", err)
		return &pb.Responsee{Message: "Failed to update goal"}, err
	}

	return &pb.Responsee{Message: "Goal updated successfully"}, nil
}

// DeleteGoal deletes a goal by its ID
func (s *GoalStorage) DeleteGoal(req *pb.DeleteGoalRequest) (*pb.GoalDeleteResponse, error) {
	coll := s.db.Collection("goals")

	_, err := coll.DeleteOne(context.Background(), bson.M{"id": req.GoalId})
	if err != nil {
		log.Printf("Failed to delete goal: %v", err)
		return &pb.GoalDeleteResponse{Success: false}, err
	}

	return &pb.GoalDeleteResponse{Success: true}, nil
}
func (s *GoalStorage) CheckGoal(ctx context.Context, userId string) (bool, string, error) {
	coll := s.db.Collection("goals")

	// Define a struct to match the document structure
	var result struct {
		TargetAmount  float32
		CurrentAmount float32
		Deadline      string
	}

	// Find the document for the given UserId
	err := coll.FindOne(ctx, bson.M{"UserId": userId}).Decode(&result)
	if err != nil {
		log.Printf("Failed to get goal by UserId: %v", err)
		return false, "", err
	}

	// Get the current date in the same string format as stored in MongoDB
	now := time.Now().Format("2006-01-02")

	// Check if 'now' matches the 'Deadline'
	if now == result.Deadline {
		if result.CurrentAmount < result.TargetAmount {
			err = s.UpdateStatusByUserId(ctx, userId, "Filed")
			if err != nil {
				log.Print("Error while update goal status")
				return false, "", err
			}
			return false, "The deadline has passed, and you did not reach your savings goal.", nil
		}
		err = s.UpdateStatusByUserId(ctx, userId, "Success")
		if err != nil {
			log.Print("Error while update goal status")
			return false, "", err
		}
		return false, "Congratulations! You reached your savings goal by the deadline.", nil
	}

	return true, "", nil
}

func (s *GoalStorage) UpdateStatusByUserId(ctx context.Context, userid string, status string) error {
	coll := s.db.Collection("goals")
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": userid}, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return err
	}
	return nil
}

func (s *GoalStorage) UpdateGoalAmount(ctx context.Context, UserId string, amount float32) error {
	coll := s.db.Collection("goals")

	update := bson.M{
		"$inc": bson.M{
			"CurrentAmount": +amount,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": UserId}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}
