package storage

import (
	"context"
	"fmt"
	"log"

	pb "budget-service/genproto"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	objID, err := primitive.ObjectIDFromHex(req.AccountId)
	if err != nil {
		return nil, fmt.Errorf("invalid budget ID: %v", err)
	}

	var accountData struct {
		ID          primitive.ObjectID `bson:"_id"`
		UserID      string             `bson:"user_id"`
		AccountName string             `bson:"account_name"`
		Type        string             `bson:"type"`
		Balance     float64            `bson:"balance"`
		Currency    string             `bson:"currency"`
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&accountData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("budget not found")
		}
		log.Printf("Failed to get account by ID: %v", err)
		return nil, err
	}

	account := &pb.AccountResponse{
		AccountId:   accountData.ID.Hex(),
		UserId:      accountData.UserID,
		AccountName: accountData.AccountName,
		AccountType: accountData.Type,
		Balance:     accountData.Balance,
		Currency:    accountData.Currency,
	}

	return account, nil
}
func (s *AccountStorage) UpdateAccount(req *pb.UpdateAccountRequest) (*pb.CreateAccountRes, error) {
	coll := s.db.Collection("accounts")

	objID, err := primitive.ObjectIDFromHex(req.AccountId)
	if err != nil {
		return &pb.CreateAccountRes{Message: "Invalid account ID"}, err
	}

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
	result, err := coll.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
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
	objID, err := primitive.ObjectIDFromHex(req.AccountId)
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, fmt.Errorf("invalid account ID: %v", err)
	}

	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
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
	// ... other filter conditions ...

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var accounts []*pb.AccountResponse
	for cursor.Next(context.Background()) {
		// Use a temporary struct for decoding
		var accountData struct {
			ID          primitive.ObjectID `bson:"_id"` // BSON tag for ID field
			UserID      string             `bson:"user_id"`
			AccountName string             `bson:"account_name"`
			Type        string             `bson:"type"`
			Balance     float64            `bson:"balance"`
			Currency    string             `bson:"currency"`
		}
		if err := cursor.Decode(&accountData); err != nil {
			log.Printf("Failed to decode account: %v", err)
			return nil, err
		}

		// Convert temporary struct data to Protobuf message
		account := &pb.AccountResponse{
			AccountId:   accountData.ID.Hex(),
			UserId:      accountData.UserID,
			AccountName: accountData.AccountName,
			AccountType: accountData.Type,
			Balance:     accountData.Balance,
			Currency:    accountData.Currency,
		}
		accounts = append(accounts, account)
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
