package storage

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "budget-service/genproto"
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

	id := uuid.NewString() // Generate a new UUID for the transaction
	_, err := coll.InsertOne(context.Background(), bson.M{
		"id":          id,
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
		var transaction pb.TransactionResponse
		if err := cursor.Decode(&transaction); err != nil {
			log.Printf("Failed to decode transaction: %v", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
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

	var transaction pb.TransactionResponse
	err := coll.FindOne(context.Background(), bson.M{"id": req.TransactionId}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No transaction found with id: %v", req.TransactionId)
			return nil, err
		}
		log.Printf("Failed to get transaction by id: %v", err)
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransaction updates a transaction based on the provided request data
func (s *TransactionStorage) UpdateTransaction(req *pb.UpdateTransactionRequest) (*pb.Response, error) {
	coll := s.db.Collection("transactions")

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

	_, err := coll.UpdateOne(context.Background(), bson.M{"id": req.TransactionId}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return &pb.Response{Message: "Failed to update transaction"}, err
	}

	return &pb.Response{Message: "Transaction updated successfully"}, nil
}

// DeleteTransaction deletes a transaction by its ID
func (s *TransactionStorage) DeleteTransaction(req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error) {
	coll := s.db.Collection("transactions")

	_, err := coll.DeleteOne(context.Background(), bson.M{"id": req.TransactionId})
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return &pb.TransactionDeleteResponse{Success: false}, err
	}

	return &pb.TransactionDeleteResponse{Success: true}, nil
}
