package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "budget-service/genproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Generate a new ObjectID for the budget
	objID := primitive.NewObjectID()
	req.Id = objID.Hex() // Set the ID field in the request

	_, err := coll.InsertOne(context.Background(), bson.M{
		"_id":         objID, // Use ObjectID for _id
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

	// Use ObjectID if budget_id is provided
	if req.BudgetId != "" {
		objID, err := primitive.ObjectIDFromHex(req.BudgetId)
		if err != nil {
			return nil, fmt.Errorf("invalid budget ID: %v", err)
		}
		filter["_id"] = objID
	}

	// Filter by start_date and end_date
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
		var budgetData struct {
			ID          primitive.ObjectID `bson:"_id"` // BSON tag for ID field
			UserID      string             `bson:"user_id"`
			CategoryID  string             `bson:"category_id"`
			Amount      float64            `bson:"amount"`
			Period      string             `bson:"period"`
			StartDate  string             `bson:"start_date"`
			EndDate    string             `bson:"end_date"`
		}
		if err := cursor.Decode(&budgetData); err != nil {
			log.Printf("Failed to decode budget: %v", err)
			return nil, err
		}

		// Convert temporary struct data to Protobuf message
		budget := &pb.BudgetResponse{
			BudgetId:    budgetData.ID.Hex(),
			UserId:      budgetData.UserID,
			CategoryId:  budgetData.CategoryID,
			Amount:      budgetData.Amount,
			Period:      budgetData.Period,
			StartDate:  budgetData.StartDate,
			EndDate:    budgetData.EndDate,
		}
		budgets = append(budgets, budget)
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

	objID, err := primitive.ObjectIDFromHex(req.BudgetId)
	if err != nil {
		return nil, fmt.Errorf("invalid budget ID: %v", err)
	}

	var budgetData struct {
		ID          primitive.ObjectID `bson:"_id"`
		UserID      string             `bson:"user_id"`
		CategoryID  string             `bson:"category_id"`
		Amount      float64            `bson:"amount"`
		Period      string             `bson:"period"`
		StartDate  string             `bson:"start_date"`
		EndDate    string             `bson:"end_date"`
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&budgetData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("budget not found")
		}
		log.Printf("Failed to get budget by ID: %v", err)
		return nil, err
	}

	budget := &pb.BudgetResponse{
		BudgetId:    budgetData.ID.Hex(),
		UserId:      budgetData.UserID,
		CategoryId:  budgetData.CategoryID,
		Amount:      budgetData.Amount,
		Period:      budgetData.Period,
		StartDate:  budgetData.StartDate,
		EndDate:    budgetData.EndDate,
	}

	return budget, nil
}

// UpdateBudget updates a budget based on the provided request data
func (s *BudgetStorage) UpdateBudget(req *pb.UpdateBudgetRequest) (*pb.MessageResponsee, error) {
	coll := s.db.Collection("budgets")

	objID, err := primitive.ObjectIDFromHex(req.BudgetId)
	if err != nil {
		return &pb.MessageResponsee{Message: "Invalid budget ID"}, err
	}

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

	_, err = coll.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update budget: %v", err)
		return &pb.MessageResponsee{Message: "Failed to update budget"}, err
	}

	return &pb.MessageResponsee{Message: "Budget updated successfully"}, nil
}

// DeleteBudget deletes a budget by its ID
func (s *BudgetStorage) DeleteBudget(req *pb.DeleteBudgetRequest) (*pb.BudgetDeleteResponse, error) {
	coll := s.db.Collection("budgets")

	objID, err := primitive.ObjectIDFromHex(req.BudgetId)
	if err != nil {
		return &pb.BudgetDeleteResponse{Success: false}, fmt.Errorf("invalid budget ID: %v", err)
	}

	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Printf("Failed to delete budget: %v", err)
		return &pb.BudgetDeleteResponse{Success: false}, err
	}

	return &pb.BudgetDeleteResponse{Success: true}, nil
}

func (s *BudgetStorage) UpdateBudgetAmount(userId string, amount float32) error {
	coll := s.db.Collection("budgets")

	update := bson.M{
		"$inc": bson.M{
			"amount": -amount,
		},
	}
	_, err := coll.UpdateOne(context.Background(), bson.M{"user_id": userId}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *BudgetStorage) CheckBudget(userId string) (bool, error) {
	coll := s.db.Collection("budgets")

	// Define a struct to match the document structure
	var result struct {
		Amount    float32  `bson:"amount"`
		StartDate string `bson:"start_date"`
		EndDate   string `bson:"end_date"`
	}

	// Find the document for the given UserId
	err := coll.FindOne(context.Background(), bson.M{"user_id": userId}).Decode(&result)
	if err != nil {
		// Other errors (e.g., database issues)
		log.Printf("Failed to get budget by UserId: %v", err)
		return false, err
	}
	// Get the current date in the same string format as stored in MongoDB
	now := time.Now().Format("2006-01-02")
	
	// Check if 'now' is between 'StartDate' and 'EndDate'
	if now >= result.StartDate && now <= result.EndDate {
		fmt.Println(result.EndDate)

		// If within the date range, check if the amount is greater than 0
		if result.Amount <= 0 {

			return false, nil
		}
	}

	return true, nil
}