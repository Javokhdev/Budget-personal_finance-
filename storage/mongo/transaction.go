package storage

import (
	"context"
	"log"

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

	_, err := coll.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: req.Id},
		{Key: "user_id", Value: req.UserId},
		{Key: "account_id", Value: req.AccountId},
		{Key: "category_id", Value: req.CategoryId},
		{Key: "amount", Value: req.Amount},
		{Key: "type", Value: req.Type},
		{Key: "description", Value: req.Description},
		{Key: "date", Value: req.Date},
	})
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return nil, err
	}

	return &pb.Response{Message: "Transaction created successfully"}, nil
}

// GetTransactions retrieves all transactions
func (s *TransactionStorage) GetTransactions(req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	coll := s.db.Collection("transactions")

	cursor, err := coll.Find(context.Background(), bson.D{})
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

	filter := bson.D{{Key: "transaction_id", Value: req.TransactionId}}
	var transaction pb.TransactionResponse
	err := coll.FindOne(context.Background(), filter).Decode(&transaction)
	if err != nil {
		log.Printf("Failed to get transaction by id: %v", err)
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransaction updates a transaction based on the provided request data
func (s *TransactionStorage) UpdateTransaction(req *pb.UpdateTransactionRequest) (*pb.Response, error) {
	coll := s.db.Collection("transactions")

	filter := bson.D{{Key: "transaction_id", Value: req.TransactionId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: req.UserId},
			{Key: "account_id", Value: req.AccountId},
			{Key: "category_id", Value: req.CategoryId},
			{Key: "amount", Value: req.Amount},
			{Key: "type", Value: req.Type},
			{Key: "description", Value: req.Description},
			{Key: "date", Value: req.Date},
		}},
	}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return nil, err
	}

	return &pb.Response{Message: "Transaction updated successfully"}, nil
}

// DeleteTransaction deletes a transaction by its ID
func (s *TransactionStorage) DeleteTransaction(req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error) {
	coll := s.db.Collection("transactions")

	filter := bson.D{{Key: "transaction_id", Value: req.TransactionId}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return nil, err
	}

	return &pb.TransactionDeleteResponse{Success: true}, nil
}
