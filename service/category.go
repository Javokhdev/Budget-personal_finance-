package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	mdb "budget-service/storage"
)

type CategoryService struct {
	stg mdb.InitRoot
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(db mdb.InitRoot) *CategoryService {
	return &CategoryService{stg: db}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.MessageResponse, error) {
	resp, err := s.stg.Category().CreateCategory(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListResponse, error) {
	resp, err := s.stg.Category().ListCategories(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error) {
	resp, err := s.stg.Category().GetCategoryById(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.MessageResponse, error) {
	resp, err := s.stg.Category().UpdateCategory(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.CategoryDeleteResponse, error) {
	resp, err := s.stg.Category().DeleteCategory(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
