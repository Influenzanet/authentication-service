package main

import (
	"log"
	"net"
	"strconv"

	auth_api "github.com/Influenzanet/api/dist/go/auth-service"
	user_api "github.com/Influenzanet/api/dist/go/user-management"
	"google.golang.org/grpc"
)

type authServiceServer struct {
}

var userManagementClient user_api.UserManagementApiClient

func connectToUserManagementServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.UserManagement, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

func init() {
	readConfig()
}

func main() {
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()

	userManagementClient = user_api.NewUserManagementApiClient(userManagementServerConn)

	log.Println("wait connections on port " + strconv.Itoa(conf.ListenPort))
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.ListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth_api.RegisterAuthServiceApiServer(grpcServer, &authServiceServer{})
	grpcServer.Serve(lis)
}
