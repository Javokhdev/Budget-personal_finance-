package storage

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "budget-service/genproto"
)

type CategoryStorage struct {
	db *mongo.Database
}

func NewCategoryStorage(db *mongo.Database) *CategoryStorage {
	return &CategoryStorage{db: db}
}

func (s *CategoryStorage) CreateCategory(req *pb.CreateCategoryRequest) (*pb.MessageResponse, error) {
	coll := s.db.Collection("categories")
	id := uuid.NewString()

	_, err := coll.InsertOne(context.Background(), bson.M{
		"id":     id,
		"user_id": req.UserId,
		"name":   req.Name,
		"type":   req.Type,
	})
	if err != nil {
		log.Printf("Failed to create category: %v", err)
		return &pb.MessageResponse{Message: "Failed to create category"}, err
	}

	return &pb.MessageResponse{Message: "Category created successfully"}, nil
}

func (s *CategoryStorage) ListCategories(req *pb.ListCategoriesRequest) (*pb.ListResponse, error) {
	coll := s.db.Collection("categories")

	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}

	cursor, err := coll.Find(context.Background(), filter)
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

func (s *CategoryStorage) GetCategoryById(req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error) {
	coll := s.db.Collection("categories")

	var category pb.CategoryResponse
	err := coll.FindOne(context.Background(), bson.M{"id": req.CategoryId}).Decode(&category)
	if err != nil {
		log.Printf("Failed to get category by id: %v", err)
		return nil, err
	}

	return &category, nil
}

func (s *CategoryStorage) UpdateCategory(req *pb.UpdateCategoryRequest) (*pb.MessageResponse, error) {
	coll := s.db.Collection("categories")

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Type != "" {
		update["type"] = req.Type
	}

	if len(update) == 0 {
		return &pb.MessageResponse{Message: "Nothing to update"}, nil
	}

	_, err := coll.UpdateOne(context.Background(), bson.M{"id": req.CategoryId}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update category: %v", err)
		return &pb.MessageResponse{Message: "Failed to update category"}, err
	}

	return &pb.MessageResponse{Message: "Category updated successfully"}, nil
}

func (s *CategoryStorage) DeleteCategory(req *pb.DeleteCategoryRequest) (*pb.CategoryDeleteResponse, error) {
	coll := s.db.Collection("categories")

	_, err := coll.DeleteOne(context.Background(), bson.M{"id": req.CategoryId})
	if err != nil {
		log.Printf("Failed to delete category: %v", err)
		return &pb.CategoryDeleteResponse{Success: false}, err
	}

	return &pb.CategoryDeleteResponse{Success: true}, nil
}
