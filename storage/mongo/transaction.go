package storage

import (
	"context"
	"fmt"
	"log"

	pb "budget-service/genproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionStorage handles transaction operations in MongoDB
type TransactionStorage struct {
	db *mongo.Database
}

// NewTransactionStorage initializes a new TransactionStorage
func NewTransactionStorage(db *mongo.Database) *TransactionStorage {
	return &TransactionStorage{db: db}
}

// CreateTransaction creates a new transaction in the database
func (s *TransactionStorage) CreateTransaction(req *pb.CreateTransactionRequest) (*pb.Response, error) {
	coll := s.db.Collection("transactions")

	// Generate a new ObjectID for the transaction
	objID := primitive.NewObjectID()
	req.Id = objID.Hex() // Set the ID field in the request

	_, err := coll.InsertOne(context.Background(), bson.M{
		"_id":          objID, // Use ObjectID for _id
		"user_id":     req.UserId,
		"account_id":  req.AccountId,
		"category_id": req.CategoryId,
		"amount":      req.Amount,
		"type":        req.Type,
		"description": req.Description,
		"date":        req.Date,
	})
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return &pb.Response{Message: "Failed to create transaction"}, err
	}

	return &pb.Response{Message: "Transaction created successfully"}, nil
}

// GetTransactions retrieves all transactions based on the filter criteria
func (s *TransactionStorage) GetTransactions(req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	coll := s.db.Collection("transactions")

	filter := bson.M{}
	if req.AccountId != "" {
		filter["account_id"] = req.AccountId
	}
	if req.CategoryId != "" {
		filter["category_id"] = req.CategoryId
	}
	if req.Amount > 0 {
		filter["amount"] = req.Amount
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}
	if req.Date != "" {
		filter["date"] = req.Date
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []*pb.TransactionResponse
	for cursor.Next(context.Background()) {
		var transactionData struct {
			ID          primitive.ObjectID `bson:"_id"` // BSON tag for ID field
			UserId      string             `bson:"user_id"`
			AccountId   string             `bson:"account_id"`
			CategoryId  string             `bson:"category_id"`
			Amount      float32            `bson:"amount"`
			Type        string             `bson:"type"`
			Description string             `bson:"description"`
			Date        string             `bson:"date"`
		}
		if err := cursor.Decode(&transactionData); err != nil {
			log.Printf("Failed to decode transaction: %v", err)
			return nil, err
		}

		transaction := &pb.TransactionResponse{
			TransactionId: transactionData.ID.Hex(),
			UserId:      transactionData.UserId,
			AccountId:   transactionData.AccountId,
			CategoryId:  transactionData.CategoryId,
			Amount:      transactionData.Amount,
			Type:        transactionData.Type,
			Description: transactionData.Description,
			Date:        transactionData.Date,
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.TransactionsResponse{Transactions: transactions}, nil
}

// GetTransactionById retrieves a transaction by its ID
func (s *TransactionStorage) GetTransactionById(req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error) {
	coll := s.db.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %v", err)
	}

	var transactionData struct {
		ID          primitive.ObjectID `bson:"_id"` // BSON tag for ID field
		UserId      string             `bson:"user_id"`
		AccountId   string             `bson:"account_id"`
		CategoryId  string             `bson:"category_id"`
		Amount      float32            `bson:"amount"`
		Type        string             `bson:"type"`
		Description string             `bson:"description"`
		Date        string             `bson:"date"`
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&transactionData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction not found")
		}
		log.Printf("Failed to get transaction by ID: %v", err)
		return nil, err
	}

	transaction := &pb.TransactionResponse{
		TransactionId: transactionData.ID.Hex(),
		UserId:      transactionData.UserId,
		AccountId:   transactionData.AccountId,
		CategoryId:  transactionData.CategoryId,
		Amount:      transactionData.Amount,
		Type:        transactionData.Type,
		Description: transactionData.Description,
		Date:        transactionData.Date,
	}

	return transaction, nil
}

// UpdateTransaction updates a transaction based on the provided request data
func (s *TransactionStorage) UpdateTransaction(req *pb.UpdateTransactionRequest) (*pb.Response, error) {
	coll := s.db.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.TransactionId)
	if err != nil {
		return &pb.Response{Message: "Invalid transaction ID"}, err
	}

	update := bson.M{}
	if req.AccountId != "" {
		update["account_id"] = req.AccountId
	}
	if req.CategoryId != "" {
		update["category_id"] = req.CategoryId
	}
	if req.Amount > 0 {
		update["amount"] = req.Amount
	}
	if req.Type != "" {
		update["type"] = req.Type
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Date != "" {
		update["date"] = req.Date
	}

	if len(update) == 0 {
		return &pb.Response{Message: "Nothing to update"}, nil
	}

	_, err = coll.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return &pb.Response{Message: "Failed to update transaction"}, err
	}

	return &pb.Response{Message: "Transaction updated successfully"}, nil
}

// DeleteTransaction deletes a transaction by its ID
func (s *TransactionStorage) DeleteTransaction(req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error) {
	coll := s.db.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.TransactionId)
	if err != nil {
		return &pb.TransactionDeleteResponse{Success: false}, fmt.Errorf("invalid transaction ID: %v", err)
	}

	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return &pb.TransactionDeleteResponse{Success: false}, err
	}

	return &pb.TransactionDeleteResponse{Success: true}, nil
}
