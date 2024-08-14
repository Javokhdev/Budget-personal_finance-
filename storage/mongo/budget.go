package storage

import (
	"context"
	"log"
	"time"

	pb "budget-service/genproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, err := coll.InsertOne(context.Background(), bson.M{
		"id":          req.Id,
		"user_id":     req.UserId,
		"category_id": req.CategoryId,
		"amount":      req.Amount,
		"period":      req.Period,
		"start_date":  req.StartDate,
		"end_date":    req.EndDate,
	})
	if err != nil {
		log.Printf("Failed to create budget: %v", err)
		return &pb.MessageResponsee{Message: "Failed to create budget"}, err
	}

	return &pb.MessageResponsee{Message: "Budget created successfully"}, nil
}

// ListBudgets lists all budgets, potentially filtering by start_date and end_date
func (s *BudgetStorage) ListBudgets(req *pb.ListBudgetsRequest) (*pb.ListBudgetsResponse, error) {
	coll := s.db.Collection("budgets")

	filter := bson.M{}

	if req.StartDate != "" && req.EndDate != "" {
		filter["start_date"] = bson.M{"$gte": req.StartDate}
		filter["end_date"] = bson.M{"$lte": req.EndDate}
	}

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

	var budget pb.BudgetResponse
	err := coll.FindOne(context.Background(), bson.M{"id": req.BudgetId}).Decode(&budget)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No budget found with id: %v", req.BudgetId)
			return nil, err
		}
		log.Printf("Failed to get budget by id: %v", err)
		return nil, err
	}

	return &budget, nil
}

// UpdateBudget updates a budget based on the provided request data
func (s *BudgetStorage) UpdateBudget(req *pb.UpdateBudgetRequest) (*pb.MessageResponsee, error) {
	coll := s.db.Collection("budgets")

	update := bson.M{}
	if req.UserId != "" {
		update["user_id"] = req.UserId
	}
	if req.CategoryId != "" {
		update["category_id"] = req.CategoryId
	}
	if req.Amount != 0 {
		update["amount"] = req.Amount
	}
	if req.Period != "" {
		update["period"] = req.Period
	}
	if req.StartDate != "" {
		update["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		update["end_date"] = req.EndDate
	}

	if len(update) == 0 {
		return &pb.MessageResponsee{Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(context.Background(), bson.M{"id": req.BudgetId}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update budget: %v", err)
		return &pb.MessageResponsee{Message: "Failed to update budget"}, err
	}

	return &pb.MessageResponsee{Message: "Budget updated successfully"}, nil
}

// DeleteBudget deletes a budget by its ID
func (s *BudgetStorage) DeleteBudget(req *pb.DeleteBudgetRequest) (*pb.BudgetDeleteResponse, error) {
	coll := s.db.Collection("budgets")

	_, err := coll.DeleteOne(context.Background(), bson.M{"id": req.BudgetId})
	if err != nil {
		log.Printf("Failed to delete budget: %v", err)
		return &pb.BudgetDeleteResponse{Success: false}, err
	}

	return &pb.BudgetDeleteResponse{Success: true}, nil
}

func (s *BudgetStorage) UpdateBudgetAmount(ctx context.Context, UserId string, amount float32) error {
	coll := s.db.Collection("budgets")

	update := bson.M{
		"$inc": bson.M{
			"amount": -amount,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"UserId": UserId}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *BudgetStorage) CheckBudget(ctx context.Context, userId string) (bool, error) {
	coll := s.db.Collection("budgets")

	// Define a struct to match the document structure
	var result struct {
		Amount    float32
		StartDate string
		EndDate   string
	}

	// Find the document for the given UserId
	err := coll.FindOne(ctx, bson.M{"UserId": userId}).Decode(&result)
	if err != nil {
		// Other errors (e.g., database issues)
		log.Printf("Failed to get budget by UserId: %v", err)
		return false, err
	}

	// Get the current date in the same string format as stored in MongoDB
	now := time.Now().Format("2006-01-02")

	// Check if 'now' is between 'StartDate' and 'EndDate'
	if now >= result.StartDate && now <= result.EndDate {

		// If within the date range, check if the amount is greater than 0
		if result.Amount <= 0 {

			return false, nil
		}
	}

	return true, nil
}
