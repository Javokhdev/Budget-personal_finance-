package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	mdb "budget-service/storage"
)

type AccountService struct {
	stg mdb.InitRoot
	pb.UnimplementedAccountServiceServer
}

func NewAccountService(db mdb.InitRoot) *AccountService {
	return &AccountService{stg: db}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountRes, error) {
	resp, err := s.stg.Account().CreateAccount(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *AccountService) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	resp, err := s.stg.Account().ListAccounts(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *AccountService) GetAccountById(ctx context.Context, req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error) {
	resp, err := s.stg.Account().GetAccountById(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.CreateAccountRes, error) {
	resp, err := s.stg.Account().UpdateAccount(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteResponse, error) {
	resp, err := s.stg.Account().DeleteAccount(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
