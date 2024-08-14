package storage

import (
	"context"
	"fmt"
	"log"

	u "budget-service/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Db    *mongo.Database
	Accounts u.AccountStorage
	Budgets u.BudgetStorage
	Categorys u.CategoryStorage
	Goals u.GoalStorage
	Transactions u.TransactionStorage
	Notifications u.NotificationService
}

func NewMongoConnection() (*MongoStorage, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", "localhost", 27017)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error: Couldn't connect to the database.", err)
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database("budget")

	return &MongoStorage{Db: db}, err
}


func (s *MongoStorage) Account() u.AccountStorage {
	if s.Accounts == nil {
		s.Accounts = &AccountStorage{s.Db}
	}
	return s.Accounts
}

func (s *MongoStorage) Budget() u.BudgetStorage {
	if s.Budgets == nil {
		s.Budgets = &BudgetStorage{s.Db}
	}
	return s.Budgets
}

func (s *MongoStorage) Category() u.CategoryStorage {
	if s.Categorys == nil {
		s.Categorys = &CategoryStorage{s.Db}
	}
	return s.Categorys
}

func (s *MongoStorage) Goal() u.GoalStorage {
	if s.Goals == nil {
		s.Goals = &GoalStorage{s.Db}
	}
	return s.Goals
}

func (s *MongoStorage) Transaction() u.TransactionStorage {
	if s.Transactions == nil {
		s.Transactions = &TransactionStorage{s.Db}
	}
	return s.Transactions
}

func (s *MongoStorage) Notification() u.NotificationService {
	if s.Notifications == nil {
		s.Notifications = &NotificationService{s.Db}
	}
	return s.Notifications
}
