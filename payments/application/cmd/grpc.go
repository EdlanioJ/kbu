package cmd

import (
	"os"

	"github.com/EdlanioJ/kbu/payments/application/config/gorm"
	"github.com/EdlanioJ/kbu/payments/application/grpc"
	"github.com/spf13/cobra"
)

var portNumber int

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		database := gorm.ConnectDB(os.Getenv("env"))

		grpc.StartGrpcServer(database, portNumber)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.Flags().IntVarP(&portNumber, "port", "p", 50051, "grpc server port")
}
