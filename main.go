package main

import (
	"log"
	"net"
	"strconv"

	api "github.com/influenzanet/authentication-service/api"
	"google.golang.org/grpc"
)

type authServiceServer struct {
}

var userManagementClient api.UserManagementApiClient

func connectToUserManagementServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.UserManagement, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

func init() {
	readConfig()
	dbInit()
}

func main() {
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()

	userManagementClient = api.NewUserManagementApiClient(userManagementServerConn)

	log.Println("wait connections on port " + strconv.Itoa(conf.ListenPort))
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.ListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceApiServer(grpcServer, &authServiceServer{})
	grpcServer.Serve(lis)
}
