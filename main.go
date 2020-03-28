package main

import (
	"log"
	"net"

	api "github.com/influenzanet/authentication-service/api"
	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc"
)

type authServiceServer struct {
}

var conf Config
var dbClient *mongo.Client
var clients = APIClients{}

// APIClients holds the service clients to the internal services
type APIClients struct {
	userManagement api.UserManagementApiClient
}

func connectToUserManagementServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.UserManagement, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

func init() {
	initConfig()
	dbInit()
	log.Println("initialization ready")
}

func main() {
	log.Println("connect to user management service")
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()
	clients.userManagement = api.NewUserManagementApiClient(userManagementServerConn)

	lis, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("wait connections on port " + conf.Port)

	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceApiServer(grpcServer, &authServiceServer{})
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
