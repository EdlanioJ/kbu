package grpc

import (
	"context"

	"github.com/EdlanioJ/kbu/payments/application/grpc/pb"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionGrpcHandler struct {
	TransactionController        *controller.Transaction
	AccountTransactionController *controller.AccountTransaction
	ServiceTransactionController *controller.ServiceTransaction
	StoreTransactionController   *controller.StoreTransaction

	pb.UnimplementedPaymentServiceServer
}

func NewTransactionGrpcHandler(
	transaction *controller.Transaction,
	accountTransaction *controller.AccountTransaction,
	serviceTransaction *controller.ServiceTransaction,
	storeTransaction *controller.StoreTransaction,
) *TransactionGrpcHandler {

	return &TransactionGrpcHandler{
		TransactionController:        transaction,
		AccountTransactionController: accountTransaction,
		ServiceTransactionController: serviceTransaction,
		StoreTransactionController:   storeTransaction,
	}
}

func (t *TransactionGrpcHandler) List(ctx context.Context, in *pb.PaginationParams) (*pb.ListResult, error) {

	result, total, err := t.TransactionController.FindAll(ctx, int(in.Page), int(in.Limit), in.Sort)

	var transactions []*pb.Transaction

	if err != nil {
		return &pb.ListResult{
			Total: 0,
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	for _, value := range result {
		transactions = append(transactions, &pb.Transaction{
			ID:        value.ID,
			Amount:    float32(value.Amount),
			Status:    value.Status,
			Currency:  value.Currency,
			FromID:    value.AccountFromID,
			ToID:      value.AccountToID,
			ServiceID: value.ServiceID,
			StoreID:   value.StoreID,
			CreatedAt: value.CreatedAt.String(),
			UpdatedAt: value.UpdatedAt.String(),
		})
	}

	return &pb.ListResult{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) Get(ctx context.Context, in *pb.Params) (*pb.PaymentResult, error) {
	result, err := t.TransactionController.Find(ctx, in.ID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ToID:      result.AccountToID,
			ServiceID: result.ServiceID,
			StoreID:   result.StoreID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) Complete(ctx context.Context, in *pb.Params) (*pb.PaymentResult, error) {
	result, err := t.TransactionController.Complete(ctx, in.ID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.Canceled, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ToID:      result.AccountToID,
			ServiceID: result.ServiceID,
			StoreID:   result.StoreID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) Error(ctx context.Context, in *pb.Params) (*pb.PaymentResult, error) {
	result, err := t.TransactionController.Error(ctx, in.ID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.Canceled, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ToID:      result.AccountToID,
			ServiceID: result.ServiceID,
			StoreID:   result.StoreID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) GetByAccountDestination(ctx context.Context, in *pb.GetByParams) (*pb.PaymentResult, error) {
	result, err := t.AccountTransactionController.FindOneByAccount(ctx, in.Id, in.TransactionID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ToID:      result.AccountToID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) ListByAccountDestination(ctx context.Context, in *pb.ListByParams) (*pb.ListResult, error) {
	result, total, err := t.AccountTransactionController.FindAllByAccountTo(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)

	var transactions []*pb.Transaction

	if err != nil {
		return &pb.ListResult{
			Total: 0,
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	for _, value := range result {
		transactions = append(transactions, &pb.Transaction{
			ID:        value.ID,
			Amount:    float32(value.Amount),
			Status:    value.Status,
			Currency:  value.Currency,
			FromID:    value.AccountFromID,
			ToID:      value.AccountToID,
			CreatedAt: value.CreatedAt.String(),
			UpdatedAt: value.UpdatedAt.String(),
		})
	}

	return &pb.ListResult{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) RegisterAccountPayment(ctx context.Context, in *pb.CreateParams) (*pb.PaymentResult, error) {
	result, err := t.AccountTransactionController.RegisterAccountTransaction(ctx, in.FromID, in.DestinationID, float64(in.Amount), in.Currency)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.Aborted, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ToID:      result.AccountToID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) GetByService(ctx context.Context, in *pb.GetByParams) (*pb.PaymentResult, error) {
	result, err := t.StoreTransactionController.FindOneByStore(ctx, in.Id, in.TransactionID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ServiceID: result.ServiceID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}
func (t *TransactionGrpcHandler) ListByService(ctx context.Context, in *pb.ListByParams) (*pb.ListResult, error) {
	result, total, err := t.ServiceTransactionController.FindAllByServiceId(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)

	var transactions []*pb.Transaction

	if err != nil {
		return &pb.ListResult{
			Total: 0,
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	for _, value := range result {
		transactions = append(transactions, &pb.Transaction{
			ID:        value.ID,
			Amount:    float32(value.Amount),
			Status:    value.Status,
			Currency:  value.Currency,
			FromID:    value.AccountFromID,
			ServiceID: value.ServiceID,
			CreatedAt: value.CreatedAt.String(),
			UpdatedAt: value.UpdatedAt.String(),
		})
	}

	return &pb.ListResult{
		Transactions: transactions,
		Total:        int32(total),
	}, nil

}
func (t *TransactionGrpcHandler) RegisterServicePayment(ctx context.Context, in *pb.CreateServiceParams) (*pb.PaymentResult, error) {
	result, err := t.ServiceTransactionController.RegisterServiceTransaction(ctx, in.FromID, in.ServiceID, in.ServicePriceID, float64(in.Amount), in.Currency)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.Aborted, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			ServiceID: result.ServiceID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) GetByStore(ctx context.Context, in *pb.GetByParams) (*pb.PaymentResult, error) {
	result, err := t.StoreTransactionController.FindOneByStore(ctx, in.Id, in.TransactionID)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			StoreID:   result.StoreID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}
func (t *TransactionGrpcHandler) ListByStore(ctx context.Context, in *pb.ListByParams) (*pb.ListResult, error) {
	result, total, err := t.StoreTransactionController.FindAllByStoreId(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)

	var transactions []*pb.Transaction

	if err != nil {
		return &pb.ListResult{
			Total: 0,
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	for _, value := range result {
		transactions = append(transactions, &pb.Transaction{
			ID:        value.ID,
			Amount:    float32(value.Amount),
			Status:    value.Status,
			Currency:  value.Currency,
			FromID:    value.AccountFromID,
			StoreID:   value.StoreID,
			CreatedAt: value.CreatedAt.String(),
			UpdatedAt: value.UpdatedAt.String(),
		})
	}

	return &pb.ListResult{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}
func (t *TransactionGrpcHandler) RegisterStorePayment(ctx context.Context, in *pb.CreateParams) (*pb.PaymentResult, error) {
	result, err := t.StoreTransactionController.RegisterStoreTransaction(ctx, in.FromID, in.DestinationID, float64(in.Amount), in.Currency)

	if err != nil {
		return &pb.PaymentResult{
			Error: err.Error(),
		}, status.Error(codes.Aborted, err.Error())
	}

	return &pb.PaymentResult{
		Transaction: &pb.Transaction{
			ID:        result.ID,
			Amount:    float32(result.Amount),
			Status:    result.Status,
			Currency:  result.Currency,
			FromID:    result.AccountFromID,
			StoreID:   result.StoreID,
			CreatedAt: result.CreatedAt.String(),
			UpdatedAt: result.UpdatedAt.String(),
		},
	}, nil
}
