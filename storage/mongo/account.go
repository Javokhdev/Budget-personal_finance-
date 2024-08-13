package storage

import (
	"context"
	"fmt"
	"log"

	pb "budget-service/genproto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStorage struct {
	db *mongo.Database
}

func (g *AccountStorage) CreateAccount(req *pb.CreateAccountRequest) (*pb.CreateAccountRes, error) {
	coll := g.db.Collection("accounts")
	_, err := coll.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: req.Id},
		{Key: "user_id", Value: req.UserId},
		{Key: "account_name", Value: req.AccountName},
		{Key: "type", Value: req.Type},
		{Key: "balance", Value: req.Balance},
		{Key: "currency", Value: req.Currency},
	})
	if err != nil {
		log.Printf("Failed to create account: %v", err)
		return nil, err
	}

	return &pb.CreateAccountRes{Message: "Account created successfully"}, nil
}

func (g *AccountStorage) ListAccounts(req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	coll := g.db.Collection("accounts")

	request := pb.AccountResponse{}

	// Construct filter conditions based on the request
	filter := bson.D{}
	if request.UserId != "" {
		filter = append(filter, bson.E{Key: "user_id", Value: request.UserId})
	}
	if request.AccountId != "" {
		filter = append(filter, bson.E{Key: "account_id", Value: request.AccountId})
	}
	if request.AccountType != "" {
		filter = append(filter, bson.E{Key: "account_type", Value: request.AccountType})
	}
	if request.Currency != "" {
		filter = append(filter, bson.E{Key: "currency", Value: request.Currency})
	}
	if request.AccountName != "" {
		filter = append(filter, bson.E{Key: "account_name", Value: request.AccountName})
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

	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}

func (g *AccountStorage) GetAccountById(req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error) {
	coll := g.db.Collection("accounts")
	var account pb.AccountResponse
	err := coll.FindOne(context.Background(), bson.D{{Key: "account_id", Value: req.AccountId}}).Decode(&account)
	if err != nil {
		log.Printf("Failed to get account by ID: %v", err)
		return nil, err
	}

	return &account, nil
}

func (g *AccountStorage) UpdateAccount(req *pb.UpdateAccountRequest) (*pb.CreateAccountRes, error) {
	// Validate input
	if req.AccountId == "" {
		return nil, fmt.Errorf("account_id is required")
	}

	coll := g.db.Collection("accounts")
	filter := bson.D{{Key: "account_id", Value: req.AccountId}}

	// Prepare the update fields
	updateFields := bson.D{}
	if req.UserId != "" {
		updateFields = append(updateFields, bson.E{Key: "user_id", Value: req.UserId})
	}
	if req.AccountName != "" {
		updateFields = append(updateFields, bson.E{Key: "account_name", Value: req.AccountName})
	}
	if req.Type != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: req.Type})
	}
	if req.Balance != 0 {
		updateFields = append(updateFields, bson.E{Key: "balance", Value: req.Balance})
	}
	if req.Currency != "" {
		updateFields = append(updateFields, bson.E{Key: "currency", Value: req.Currency})
	}

	// If no fields to update, return an appropriate message
	if len(updateFields) == 0 {
		return &pb.CreateAccountRes{Message: "No fields to update"}, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	// Perform the update
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return nil, err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return &pb.CreateAccountRes{Message: "No account found with the provided account_id"}, nil
	}

	return &pb.CreateAccountRes{Message: "Account updated successfully"}, nil
}

func (g *AccountStorage) DeleteAccount(req *pb.DeleteAccountRequest) (*pb.DeleteResponse, error) {
	coll := g.db.Collection("accounts")
	_, err := coll.DeleteOne(context.Background(), bson.D{{Key: "account_id", Value: req.AccountId}})
	if err != nil {
		log.Printf("Failed to delete account: %v", err)
		return nil, err
	}

	return &pb.DeleteResponse{Success: true}, nil
}
