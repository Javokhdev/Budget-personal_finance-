package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "budget-service/genproto"
)

// CategoryStorage struct to handle category operations in MongoDB
type CategoryStorage struct {
	db *mongo.Database
}

// NewCategoryStorage initializes a new CategoryStorage
func NewCategoryStorage(db *mongo.Database) *CategoryStorage {
	return &CategoryStorage{db: db}
}

// CreateCategory creates a new category in the database
func (s *CategoryStorage) CreateCategory(req *pb.CreateCategoryRequest) (*pb.MessageResponse, error) {
	coll := s.db.Collection("categories")

	_, err := coll.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: req.Id},
		{Key: "user_id", Value: req.UserId},
		{Key: "name", Value: req.Name},
		{Key: "type", Value: req.Type},
	})
	if err != nil {
		log.Printf("Failed to create category: %v", err)
		return nil, err
	}

	return &pb.MessageResponse{Message: "Category created successfully"}, nil
}

// ListCategories lists all categories
func (s *CategoryStorage) ListCategories(req *pb.ListCategoriesRequest) (*pb.ListResponse, error) {
	coll := s.db.Collection("categories")

	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Failed to list categories: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var categories []*pb.CategoryResponse
	for cursor.Next(context.Background()) {
		var category pb.CategoryResponse
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Failed to decode category: %v", err)
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return &pb.ListResponse{Categories: categories}, nil
}

// GetCategoryById retrieves a category by its ID
func (s *CategoryStorage) GetCategoryById(req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error) {
	coll := s.db.Collection("categories")

	filter := bson.D{{Key: "id", Value: req.CategoryId}}
	var category pb.CategoryResponse
	err := coll.FindOne(context.Background(), filter).Decode(&category)
	if err != nil {
		log.Printf("Failed to get category by id: %v", err)
		return nil, err
	}

	return &category, nil
}

// UpdateCategory updates a category based on the provided request data
func (s *CategoryStorage) UpdateCategory(req *pb.UpdateCategoryRequest) (*pb.MessageResponse, error) {
	coll := s.db.Collection("categories")

	filter := bson.D{{Key: "id", Value: req.CategoryId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: req.UserId},
			{Key: "name", Value: req.Name},
			{Key: "type", Value: req.Type},
		}},
	}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update category: %v", err)
		return nil, err
	}

	return &pb.MessageResponse{Message: "Category updated successfully"}, nil
}

// DeleteCategory deletes a category by its ID
func (s *CategoryStorage) DeleteCategory(req *pb.DeleteCategoryRequest) (*pb.CategoryDeleteResponse, error) {
	coll := s.db.Collection("categories")

	filter := bson.D{{Key: "id", Value: req.CategoryId}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Failed to delete category: %v", err)
		return nil, err
	}

	return &pb.CategoryDeleteResponse{Success: true}, nil
}
