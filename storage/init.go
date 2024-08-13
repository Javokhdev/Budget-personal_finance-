package storage

import (
	pb "budget-service/genproto"
)

type InitRoot interface {
	Account() AccountStorage
	Budget() BudgetStorage
	Category() CategoryStorage
	Goal() GoalStorage
	Transaction() TransactionStorage 
}

type AccountStorage interface {
	CreateAccount(req *pb.CreateAccountRequest)  (*pb.CreateAccountRes, error)
	ListAccounts(req *pb.ListAccountsRequest)  (*pb.ListAccountsResponse, error)
	GetAccountById(req *pb.GetAccountByIdRequest)  (*pb.AccountResponse, error)
	UpdateAccount(req *pb.UpdateAccountRequest)  (*pb.CreateAccountRes, error)
	DeleteAccount(req *pb.DeleteAccountRequest)  (*pb.DeleteResponse, error)
}

type BudgetStorage interface {
	CreateBudget(req *pb.CreateBudgetRequest)  (*pb.MessageResponsee, error)
    ListBudgets(req *pb.ListBudgetsRequest)  (*pb.ListBudgetsResponse, error)
    GetBudgetById(req *pb.GetBudgetByIdRequest)  (*pb.BudgetResponse, error)
    UpdateBudget(req *pb.UpdateBudgetRequest)  (*pb.MessageResponsee, error)
    DeleteBudget(req *pb.DeleteBudgetRequest)  (*pb.BudgetDeleteResponse, error)
}

type CategoryStorage interface {
	CreateCategory(req *pb.CreateCategoryRequest)  (*pb.MessageResponse, error)
	ListCategories(req *pb.ListCategoriesRequest)  (*pb.ListResponse, error)
	GetCategoryById(req *pb.GetCategoryByIdRequest)  (*pb.CategoryResponse, error)
	UpdateCategory(req *pb.UpdateCategoryRequest)  (*pb.MessageResponse, error)
	DeleteCategory(req *pb.DeleteCategoryRequest)  (*pb.CategoryDeleteResponse, error)
}

type GoalStorage interface {
	CreateGoal(req *pb.CreateGoalRequest)  (*pb.Responsee, error)
	ListGoals(req *pb.ListGoalsRequest)  (*pb.ListGoalsResponse, error)
	GetGoalById(req *pb.GetGoalByIdRequest)  (*pb.GoalResponse, error)
	UpdateGoal(req *pb.UpdateGoalRequest)  (*pb.Responsee, error)
	DeleteGoal(req *pb.DeleteGoalRequest)  (*pb.GoalDeleteResponse, error)
}

type TransactionStorage interface {
    CreateTransaction(req *pb.CreateTransactionRequest) (*pb.Response, error)
    GetTransactions(req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error)
    GetTransactionById(req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error)
    UpdateTransaction(req *pb.UpdateTransactionRequest) (*pb.Response, error)
    DeleteTransaction(req *pb.DeleteTransactionRequest) (*pb.TransactionDeleteResponse, error)
}
