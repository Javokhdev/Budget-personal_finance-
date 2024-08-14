package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	"budget-service/kafka"
	"budget-service/model"
	mdb "budget-service/storage"
)

func ConnectToKafka() kafka.KafkaProducer {
	kaf, err := kafka.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatal("Error while connection kafka: ", err.Error())
	}
	return kaf

}

type TransactionService struct {
	stg mdb.InitRoot
	pb.UnimplementedTransactionServiceServer
}

func NewTransactionService(db mdb.InitRoot) *TransactionService {
	return &TransactionService{stg: db}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.Response, error) {
	// Create the transaction
	resp, err := s.stg.Transaction().CreateTransaction(req)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return &pb.Response{Message: "Failed to create transaction"}, err
	}

	// Handle balance updates and notifications
	if req.Type == "-" {
		// Withdraw: Update account balance and budget amount
		err = s.stg.Account().UpdateBalanceMinus(ctx, req.AccountId, req.Amount)
		if err != nil {
			log.Printf("Failed to update account balance: %v", err)
			return &pb.Response{Message: "Failed to update account balance"}, err
		}

		err = s.stg.Budget().UpdateBudgetAmount(ctx, req.UserId, req.Amount)
		if err != nil {
			log.Printf("Failed to update budget amount: %v", err)
			return &pb.Response{Message: "Failed to update budget amount"}, err
		}

		check, err := s.stg.Budget().CheckBudget(ctx, req.UserId)
		if err != nil {
			log.Printf("Failed to check budget: %v", err)
			return &pb.Response{Message: "Failed to check budget"}, err
		}
		if !check {
			kafkaConn := ConnectToKafka()
			request := model.Send{Message: "Your Budget is depleted", UserId: req.UserId}
			err = kafka.CreateNotification(kafkaConn, &request)
			if err != nil {
				log.Printf("Failed to send Kafka notification: %v", err)
				return &pb.Response{Message: "Failed to send notification"}, err
			}
		}
	} else {
		// Deposit: Update account balance and goal amount
		err = s.stg.Account().UpdateBalance(ctx, req.AccountId, req.Amount)
		if err != nil {
			log.Printf("Failed to update account balance: %v", err)
			return &pb.Response{Message: "Failed to update account balance"}, err
		}

		err = s.stg.Goal().UpdateGoalAmount(ctx, req.UserId, req.Amount)
		if err != nil {
			log.Printf("Failed to update goal amount: %v", err)
			return &pb.Response{Message: "Failed to update goal amount"}, err
		}

		goalCheck, message, err := s.stg.Goal().CheckGoal(ctx, req.UserId)
		if err != nil {
			log.Printf("Failed to check goal: %v", err)
			return &pb.Response{Message: "Failed to check goal"}, err
		}
		if !goalCheck {
			kafkaConn := ConnectToKafka()
			request := model.Send{Message: message, UserId: req.UserId}
			err = kafka.CreateNotification(kafkaConn, &request)
			if err != nil {
				log.Printf("Failed to send Kafka notification: %v", err)
				return &pb.Response{Message: "Failed to send notification"}, err
			}
		}
	}

	return resp, nil
}

func (s *TransactionService) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	resp, err := s.stg.Transaction().GetTransactions(req)
	if err != nil {
		log.Printf("Failed to get transactions: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error) {
	resp, err := s.stg.Transaction().GetTransactionById(req)
	if err != nil {
		log.Printf("Failed to get transaction by ID: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.Response, error) {
	resp, err := s.stg.Transaction().UpdateTransaction(req)
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		return &pb.Response{Message: "Failed to update transaction"}, err
	}
	return resp, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error) {
	resp, err := s.stg.Transaction().DeleteTransaction(req)
	if err != nil {
		log.Printf("Failed to delete transaction: %v", err)
		return nil, err
	}
	return resp, nil
}
