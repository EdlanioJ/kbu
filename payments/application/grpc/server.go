package grpc

import (
	"fmt"
	"net"

	"github.com/EdlanioJ/kbu/payments/application/factory"
	"github.com/EdlanioJ/kbu/payments/application/grpc/pb"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	transactionController := factory.TransactionControllerFactory(database)

	grpcHandler := NewTransactionGrpcHandler(
		transactionController,
	)

	pb.RegisterPaymentServiceServer(grpcServer, grpcHandler)

	address := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Error(err)
	}

	log.Infof("gRPC server has been started on port %d", port)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}
}
