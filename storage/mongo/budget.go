package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "budget-service/genproto"
)

// BudgetStorage struct to handle budget operations in MongoDB
type BudgetStorage struct {
	db *mongo.Database
}

// NewBudgetStorage initializes a new BudgetStorage
func NewBudgetStorage(db *mongo.Database) *BudgetStorage {
	return &BudgetStorage{db: db}
}

// CreateBudget creates a new budget in the database
func (s *BudgetStorage) CreateBudget(req *pb.CreateBudgetRequest) (*pb.MessageResponsee, error) {
	coll := s.db.Collection("budgets")

	_, err := coll.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: req.Id},
		{Key: "user_id", Value: req.UserId},
		{Key: "category_id", Value: req.CategoryId},
		{Key: "amount", Value: req.Amount},
		{Key: "period", Value: req.Period},
		{Key: "start_date", Value: req.StartDate},
		{Key: "end_date", Value: req.EndDate},
	})
	if err != nil {
		log.Printf("Failed to create budget: %v", err)
		return nil, err
	}

	return &pb.MessageResponsee{Message: "Budget created successfully"}, nil
}

// ListBudgets lists all budgets, potentially filtering by start_date and end_date
func (s *BudgetStorage) ListBudgets(req *pb.ListBudgetsRequest) (*pb.ListBudgetsResponse, error) {
	coll := s.db.Collection("budgets")

	filter := bson.D{}

	// if req.StartDate != "" && req.EndDate != "" {
	// 	filter = append(filter, bson.E{Key: "start_date", Value: bson.M{"$gte": req.StartDate}})
	// 	filter = append(filter, bson.E{Key: "end_date", Value: bson.M{"$lte": req.EndDate}})
	// }

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to list budgets: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var budgets []*pb.BudgetResponse
	for cursor.Next(context.Background()) {
		var budget pb.BudgetResponse
		if err := cursor.Decode(&budget); err != nil {
			log.Printf("Failed to decode budget: %v", err)
			return nil, err
		}
		budgets = append(budgets, &budget)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListBudgetsResponse{Budgets: budgets}, nil
}

// GetBudgetById retrieves a budget by its ID
func (s *BudgetStorage) GetBudgetById(req *pb.GetBudgetByIdRequest) (*pb.BudgetResponse, error) {
	coll := s.db.Collection("budgets")

	filter := bson.D{{Key: "id", Value: req.BudgetId}}
	var budget pb.BudgetResponse
	err := coll.FindOne(context.Background(), filter).Decode(&budget)
	if err != nil {
		log.Printf("Failed to get budget by id: %v", err)
		return nil, err
	}

	return &budget, nil
}

// UpdateBudget updates a budget based on the provided request data
func (s *BudgetStorage) UpdateBudget(req *pb.UpdateBudgetRequest) (*pb.MessageResponsee, error) {
	coll := s.db.Collection("budgets")

	filter := bson.D{{Key: "id", Value: req.BudgetId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: req.UserId},
			{Key: "category_id", Value: req.CategoryId},
			{Key: "amount", Value: req.Amount},
			{Key: "period", Value: req.Period},
			{Key: "start_date", Value: req.StartDate},
			{Key: "end_date", Value: req.EndDate},
		}},
	}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update budget: %v", err)
		return nil, err
	}

	return &pb.MessageResponsee{Message: "Budget updated successfully"}, nil
}

// DeleteBudget deletes a budget by its ID
func (s *BudgetStorage) DeleteBudget(req *pb.DeleteBudgetRequest) (*pb.BudgetDeleteResponse, error) {
	coll := s.db.Collection("budgets")

	filter := bson.D{{Key: "id", Value: req.BudgetId}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to delete budget: %v", err)
		return nil, err
	}

	return &pb.BudgetDeleteResponse{Success: true}, nil
}
