package main

import (
	"log"
	"net"
	pb "budget-service/genproto"
	"budget-service/service"
	postgres "budget-service/storage/mongo"
	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.NewMongoConnection()
	if err != nil {
		log.Fatal("Error while connection on db: ", err.Error())
	}
	liss, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal("Error while connection on tcp: ", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, service.NewAccountService(db))
	pb.RegisterCategoryServiceServer(s, service.NewCategoryService(db))
	pb.RegisterTransactionServiceServer(s, service.NewTransactionService(db))
	pb.RegisterGoalServiceServer(s, service.NewGoalService(db))
	pb.RegisterBudgetServiceServer(s, service.NewBudgetService(db))
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}
