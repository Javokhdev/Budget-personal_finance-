package service

import (
	"context"
	"log"

	pb "budget-service/genproto"
	mdb "budget-service/storage"
)

type GoalService struct {
	stg mdb.InitRoot
	pb.UnimplementedGoalServiceServer
}

func NewGoalService(db mdb.InitRoot) *GoalService {
	return &GoalService{stg: db}
}

func (s *GoalService) CreateGoal(ctx context.Context, req *pb.CreateGoalRequest) (*pb.Responsee, error) {
	resp, err := s.stg.Goal().CreateGoal(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *GoalService) ListGoals(ctx context.Context, req *pb.ListGoalsRequest) (*pb.ListGoalsResponse, error) {
	resp, err := s.stg.Goal().ListGoals(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *GoalService) GetGoalById(ctx context.Context, req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error) {
	resp, err := s.stg.Goal().GetGoalById(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *GoalService) UpdateGoal(ctx context.Context, req *pb.UpdateGoalRequest) (*pb.Responsee, error) {
	resp, err := s.stg.Goal().UpdateGoal(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}

func (s *GoalService) DeleteGoal(ctx context.Context, req *pb.DeleteGoalRequest) (*pb.GoalDeleteResponse, error) {
	resp, err := s.stg.Goal().DeleteGoal(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
