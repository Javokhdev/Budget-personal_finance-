package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	mdb "budget-service/storage"
)

type TransactionService struct {
	stg mdb.InitRoot
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionService(db mdb.InitRoot) *TransactionService {
	return &TransactionService{stg: db}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.Response, error) {
	resp, err := s.stg.Transaction().CreateTransaction(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	resp, err := s.stg.Transaction().GetTransactions(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error) {
	resp, err := s.stg.Transaction().GetTransactionById(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.Response, error) {
	resp, err := s.stg.Transaction().UpdateTransaction(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error) {
	resp, err := s.stg.Transaction().DeleteTransaction(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
