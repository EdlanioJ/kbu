package grpc

import (
	"context"

	"github.com/EdlanioJ/kbu/payments/application/grpc/pb"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionGrpcHandler struct {
	TransactionController *controller.Transaction

	pb.UnimplementedPaymentServiceServer
}

func NewTransactionGrpcHandler(
	transaction *controller.Transaction,
) *TransactionGrpcHandler {

	return &TransactionGrpcHandler{
		TransactionController: transaction,
	}
}

func (t *TransactionGrpcHandler) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.Response, error) {
	response, err := t.TransactionController.Register(ctx, in.AccountFrom, in.AccountFrom, in.ExternalID, in.Type.String(), in.Currency, float64(in.Amount))

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}
func (t *TransactionGrpcHandler) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	response, err := t.TransactionController.Get(ctx, in.ID)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) List(ctx context.Context, in *pb.PaginationRequest) (*pb.ListResponse, error) {
	response, total, err := t.TransactionController.List(ctx, int(in.Page), int(in.Limit), in.Sort)
	var transactions []*pb.Transaction

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if transactions == nil {
		return nil, status.Error(codes.NotFound, "no payment was found")
	}

	for _, value := range response {
		transactions = append(transactions, &pb.Transaction{
			ID:          value.ID,
			Amount:      float32(value.Amount),
			Status:      value.Status,
			Currency:    value.Currency,
			AccountFrom: value.AccountFromID,
			AccountTo:   value.AccountToID,
			Type:        value.Type,
			ExternalID:  value.ExternalID,
			CreatedAt:   value.CreatedAt.String(),
			UpdatedAt:   value.UpdatedAt.String(),
		})
	}

	return &pb.ListResponse{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) GetByType(ctx context.Context, in *pb.GetByTypeRequest) (*pb.Response, error) {
	response, err := t.TransactionController.GetByType(ctx, in.TransactionID, in.Type.String())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) ListByType(ctx context.Context, in *pb.ListByTypeRequest) (*pb.ListResponse, error) {
	response, total, err := t.TransactionController.ListByType(ctx, in.Type.String(), int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)
	var transactions []*pb.Transaction

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if transactions == nil {
		return nil, status.Error(codes.NotFound, "no payment was found")
	}

	for _, value := range response {
		transactions = append(transactions, &pb.Transaction{
			ID:          value.ID,
			Amount:      float32(value.Amount),
			Status:      value.Status,
			Currency:    value.Currency,
			AccountFrom: value.AccountFromID,
			AccountTo:   value.AccountToID,
			Type:        value.Type,
			ExternalID:  value.ExternalID,
			CreatedAt:   value.CreatedAt.String(),
			UpdatedAt:   value.UpdatedAt.String(),
		})
	}

	return &pb.ListResponse{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) GetByReference(ctx context.Context, in *pb.GetRequest) (*pb.Response, error) {
	response, err := t.TransactionController.GetByExternalID(ctx, in.TransactionID, in.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) ListByReference(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	response, total, err := t.TransactionController.ListByExternalID(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)
	var transactions []*pb.Transaction

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if transactions == nil {
		return nil, status.Error(codes.NotFound, "no payment was found")
	}

	for _, value := range response {
		transactions = append(transactions, &pb.Transaction{
			ID:          value.ID,
			Amount:      float32(value.Amount),
			Status:      value.Status,
			Currency:    value.Currency,
			AccountFrom: value.AccountFromID,
			AccountTo:   value.AccountToID,
			Type:        value.Type,
			ExternalID:  value.ExternalID,
			CreatedAt:   value.CreatedAt.String(),
			UpdatedAt:   value.UpdatedAt.String(),
		})
	}

	return &pb.ListResponse{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) GetByAccountFrom(ctx context.Context, in *pb.GetRequest) (*pb.Response, error) {
	response, err := t.TransactionController.GetByAccountFrom(ctx, in.TransactionID, in.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) ListByAccountFrom(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	response, total, err := t.TransactionController.ListByAccountFrom(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)
	var transactions []*pb.Transaction

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if transactions == nil {
		return nil, status.Error(codes.NotFound, "no payment was found")
	}

	for _, value := range response {
		transactions = append(transactions, &pb.Transaction{
			ID:          value.ID,
			Amount:      float32(value.Amount),
			Status:      value.Status,
			Currency:    value.Currency,
			AccountFrom: value.AccountFromID,
			AccountTo:   value.AccountToID,
			Type:        value.Type,
			ExternalID:  value.ExternalID,
			CreatedAt:   value.CreatedAt.String(),
			UpdatedAt:   value.UpdatedAt.String(),
		})
	}

	return &pb.ListResponse{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}

func (t *TransactionGrpcHandler) GetByAccountTo(ctx context.Context, in *pb.GetRequest) (*pb.Response, error) {
	response, err := t.TransactionController.GetByAccoutTo(ctx, in.TransactionID, in.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Response{
		Transaction: &pb.Transaction{
			ID:          response.ID,
			Amount:      float32(response.Amount),
			Status:      response.Status,
			Currency:    response.Currency,
			AccountFrom: response.AccountFromID,
			AccountTo:   response.AccountToID,
			Type:        response.Type,
			ExternalID:  response.ExternalID,
			CreatedAt:   response.CreatedAt.String(),
			UpdatedAt:   response.UpdatedAt.String(),
		},
	}, nil
}

func (t *TransactionGrpcHandler) ListByAccountTo(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	response, total, err := t.TransactionController.ListByAccountTo(ctx, in.ID, int(in.Pagination.Page), int(in.Pagination.Limit), in.Pagination.Sort)
	var transactions []*pb.Transaction

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if transactions == nil {
		return nil, status.Error(codes.NotFound, "no payment was found")
	}

	for _, value := range response {
		transactions = append(transactions, &pb.Transaction{
			ID:          value.ID,
			Amount:      float32(value.Amount),
			Status:      value.Status,
			Currency:    value.Currency,
			AccountFrom: value.AccountFromID,
			AccountTo:   value.AccountToID,
			Type:        value.Type,
			ExternalID:  value.ExternalID,
			CreatedAt:   value.CreatedAt.String(),
			UpdatedAt:   value.UpdatedAt.String(),
		})
	}

	return &pb.ListResponse{
		Transactions: transactions,
		Total:        int32(total),
	}, nil
}
