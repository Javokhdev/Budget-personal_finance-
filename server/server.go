package main

import (
	pb "budget-service/genproto"
	"budget-service/service"
	postgres "budget-service/storage/mongo"
	"log"
	"net"
	// "budget-service/kafka"
	// kaf "budget-service/notificationKafka"
	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.NewMongoConnection()
	if err != nil {
		log.Fatal("Error while connection on db: ", err.Error())
	}
	// kcm := kafka.NewKafkaConsumerManager()
	// appService := service.NewNotificationService(db)
	// brokers := []string{"localhost:9092"}
	// if err := kcm.RegisterConsumer(brokers, "create", "root", kaf.StartLevel(appService)); err != nil {
	// 	if err == kafka.ErrConsumerAlreadyExists {
	// 		log.Printf("Consumer for topic 'create-job_application' already exists")
	// 	} else {
	// 		log.Fatalf("Error registering consumer: %v", err)
	// 	}
	// }

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
	pb.RegisterNotificationtServiceServer(s, service.NewNotificationService(db))
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
