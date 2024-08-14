package storage

import (
	"context"
	"fmt"
	"log"

	pb "budget-service/genproto"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStorage struct {
	db *mongo.Database
}

func NewAccountStorage(db *mongo.Database) *AccountStorage {
	return &AccountStorage{db: db}
}

func (s *AccountStorage) CreateAccount(req *pb.CreateAccountRequest) (*pb.CreateAccountRes, error) {
	coll := s.db.Collection("accounts")
	id := uuid.NewString()
	_, err := coll.InsertOne(context.Background(), bson.M{
		"id":           id,
		"user_id":      req.UserId,
		"account_name": req.AccountName,
		"type":         req.Type,
		"balance":      req.Balance,
		"currency":     req.Currency,
	})
	if err != nil {
		log.Printf("Failed to create account: %v", err)
		return &pb.CreateAccountRes{
			Message: "Failed to create account",
		}, err
	}

	return &pb.CreateAccountRes{Message: "Account created successfully"}, nil
}

func (s *AccountStorage) GetAccountById(req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error) {
	coll := s.db.Collection("accounts")
	var account pb.AccountResponse
	err := coll.FindOne(context.Background(), bson.M{"id": req.AccountId}).Decode(&account)
	if err != nil {
		log.Printf("Failed to get account by ID: %v", err)
		return nil, err
	}

	return &account, nil
}

func (s *AccountStorage) UpdateAccount(req *pb.UpdateAccountRequest) (*pb.CreateAccountRes, error) {
	coll := s.db.Collection("accounts")

	// Prepare the update fields
	updateFields := bson.M{}
	if req.UserId != "" {
		updateFields["user_id"] = req.UserId
	}
	if req.AccountName != "" {
		updateFields["account_name"] = req.AccountName
	}
	if req.Type != "" {
		updateFields["type"] = req.Type
	}
	if req.Balance != 0 {
		updateFields["balance"] = req.Balance
	}
	if req.Currency != "" {
		updateFields["currency"] = req.Currency
	}

	// If no fields to update, return an appropriate message
	if len(updateFields) == 0 {
		return &pb.CreateAccountRes{Message: "No fields to update"}, nil
	}

	update := bson.M{"$set": updateFields}
	result, err := coll.UpdateOne(context.Background(), bson.M{"id": req.AccountId}, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return &pb.CreateAccountRes{
			Message: "Failed to update account",
		}, err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return &pb.CreateAccountRes{Message: "No account found with the provided account_id"}, nil
	}

	return &pb.CreateAccountRes{Message: "Account updated successfully"}, nil
}

func (s *AccountStorage) DeleteAccount(req *pb.DeleteAccountRequest) (*pb.DeleteResponse, error) {
	coll := s.db.Collection("accounts")
	_, err := coll.DeleteOne(context.Background(), bson.M{"id": req.AccountId})
	if err != nil {
		log.Printf("Failed to delete account: %v", err)
		return &pb.DeleteResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteResponse{Success: true}, nil
}

func (s *AccountStorage) ListAccounts(req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	coll := s.db.Collection("accounts")
	filter := bson.M{}

	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.AccountId != "" {
		filter["id"] = req.AccountId
	}
	if req.AccountType != "" {
		filter["type"] = req.AccountType
	}
	if req.Currency != "" {
		filter["currency"] = req.Currency
	}
	if req.AccountName != "" {
		filter["account_name"] = req.AccountName
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var accounts []*pb.AccountResponse
	for cursor.Next(context.Background()) {
		var account pb.AccountResponse
		if err := cursor.Decode(&account); err != nil {
			log.Printf("Failed to decode account: %v", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}

func (s *AccountStorage) UpdateBalance(ctx context.Context, accountID string, amount float32) error {
	coll := s.db.Collection("accounts")

	// Use the $inc operator to add the amount to the existing balance
	update := bson.M{
		"$inc": bson.M{
			"balance": amount,
		},
	}

	_, err := coll.UpdateOne(ctx, bson.M{"id": accountID}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}
	return nil
}

func (s *AccountStorage) UpdateBalanceMinus(ctx context.Context, accountID string, amount float32) error {
	coll := s.db.Collection("accounts")

	// Use the $inc operator to decrement the balance by the given amount
	update := bson.M{
		"$inc": bson.M{
			"balance": -amount,
		},
	}

	// Perform the update operation
	result, err := coll.UpdateOne(ctx, bson.M{"id": accountID}, update)
	if err != nil {
		log.Printf("Failed to update account balance: %v", err)
		return err
	}

	// Check if any document was matched by the query
	if result.MatchedCount == 0 {
		err = fmt.Errorf("no account found with ID %s", accountID)
		log.Printf("Failed to update account balance: %v", err)
		return err
	}

	return nil
}
