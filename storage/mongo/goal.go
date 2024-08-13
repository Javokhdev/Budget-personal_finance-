package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "budget-service/genproto"
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

	_, err := coll.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: req.Id},
		{Key: "user_id", Value: req.UserId},
		{Key: "target_amount", Value: req.TargetAmount},
		{Key: "current_amount", Value: req.CurrentAmount},
		{Key: "deadline", Value: req.Deadline},
		{Key: "status", Value: req.Status},
	})
	if err != nil {
		log.Printf("Failed to create goal: %v", err)
		return nil, err
	}

	return &pb.Responsee{Message: "Goal created successfully"}, nil
}

// ListGoals lists all goals
func (s *GoalStorage) ListGoals(req *pb.ListGoalsRequest) (*pb.ListGoalsResponse, error) {
	coll := s.db.Collection("goals")

	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Failed to list goals: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var goals []*pb.GoalResponse
	for cursor.Next(context.Background()) {
		var goal pb.GoalResponse
		if err := cursor.Decode(&goal); err != nil {
			log.Printf("Failed to decode goal: %v", err)
			return nil, err
		}
		goals = append(goals, &goal)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListGoalsResponse{Goals: goals}, nil
}

// GetGoalById retrieves a goal by its ID
func (s *GoalStorage) GetGoalById(req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error) {
	coll := s.db.Collection("goals")

	filter := bson.D{{Key: "id", Value: req.GoalId}}
	var goal pb.GoalResponse
	err := coll.FindOne(context.Background(), filter).Decode(&goal)
	if err != nil {
		log.Printf("Failed to get goal by id: %v", err)
		return nil, err
	}

	return &goal, nil
}

// UpdateGoal updates a goal based on the provided request data
func (s *GoalStorage) UpdateGoal(req *pb.UpdateGoalRequest) (*pb.Responsee, error) {
	coll := s.db.Collection("goals")

	filter := bson.D{{Key: "id", Value: req.GoalId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: req.UserId},
			{Key: "target_amount", Value: req.TargetAmount},
			{Key: "current_amount", Value: req.CourentAmount},
			{Key: "deadline", Value: req.Deadline},
		}},
	}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update goal: %v", err)
		return nil, err
	}

	return &pb.Responsee{Message: "Goal updated successfully"}, nil
}

// DeleteGoal deletes a goal by its ID
func (s *GoalStorage) DeleteGoal(req *pb.DeleteGoalRequest) (*pb.GoalDeleteResponse, error) {
	coll := s.db.Collection("goals")

	filter := bson.D{{Key: "id", Value: req.GoalId}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to delete goal: %v", err)
		return nil, err
	}

	return &pb.GoalDeleteResponse{Success: true}, nil
}
