package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	mdb "budget-service/storage"
)

type BudgetService struct {
	stg mdb.InitRoot
	pb.UnimplementedBudgetServiceServer
}

func NewBudgetService(db mdb.InitRoot) *BudgetService {
	return &BudgetService{stg: db}
}

func (s *BudgetService) CreateBudget(ctx context.Context, req *pb.CreateBudgetRequest) (*pb.MessageResponsee, error) {
	resp, err := s.stg.Budget().CreateBudget(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *BudgetService) ListBudgets(ctx context.Context, req *pb.ListBudgetsRequest) (*pb.ListBudgetsResponse, error) {
	resp, err := s.stg.Budget().ListBudgets(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *BudgetService) GetBudgetById(ctx context.Context, req *pb.GetBudgetByIdRequest) (*pb.BudgetResponse, error) {
	resp, err := s.stg.Budget().GetBudgetById(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *BudgetService) UpdateBudget(ctx context.Context, req *pb.UpdateBudgetRequest) (*pb.MessageResponsee, error) {
	resp, err := s.stg.Budget().UpdateBudget(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *BudgetService) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetRequest) (*pb.BudgetDeleteResponse, error) {
	resp, err := s.stg.Budget().DeleteBudget(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
