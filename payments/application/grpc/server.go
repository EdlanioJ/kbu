package grpc

import (
	"fmt"
	"net"

	"github.com/EdlanioJ/kbu/payments/application/factory"
	"github.com/EdlanioJ/kbu/payments/application/grpc/pb"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	transactionController := factory.TransactionControllerFactory(database)
	serviceTransactionController := factory.ServiceTransactionControllerFactory(database)
	accountTransactionController := factory.AccountTransactionControllerFactory(database)
	storeTransactionController := factory.StoreTransactionControllerFactory(database)

	grpcHandler := NewTransactionGrpcHandler(
		transactionController,
		accountTransactionController,
		serviceTransactionController,
		storeTransactionController,
	)

	pb.RegisterPaymentServiceServer(grpcServer, grpcHandler)

	address := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msgf("gRPC server has been started on port %d", port)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
