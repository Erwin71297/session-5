package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"session-5/common/config"
	"session-5/common/models"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoServer struct {
	models.UnimplementedToDoServiceServer
}

var Todos = make(map[string]*models.ToDo)

func main() {
	srv := grpc.NewServer()
	userSrv := new(TodoServer)
	models.RegisterToDoServiceServer(srv, userSrv)

	log.Println("Starting Todo Server at ", config.SERVICE_TODO_PORT)

	listener, err := net.Listen("tcp", config.SERVICE_TODO_PORT)
	if err != nil {
		log.Fatalf("could not listen. Err: %+v\n", err)
	}

	// setup http proxy
	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost"+config.SERVICE_TODO_PORT, "gRPC server endpoint")
		_ = models.RegisterToDoServiceHandlerFromEndpoint(context.Background(), mux, *grpcServerEndpoint, opts)
		log.Println("Starting Todo Server HTTP at 9001 ")
		http.ListenAndServe(":9001", mux)
	}()

	log.Fatalln(srv.Serve(listener))
}

func (t *TodoServer) ReadAll(ctx context.Context, req *emptypb.Empty) (*models.ReadAllResponse, error) {
	var todos []*models.ToDo
	for _, v := range Todos {
		todos = append(todos, &models.ToDo{
			ID:   v.GetID(),
			Name: v.GetName(),
		})
	}
	return &models.ReadAllResponse{Todo: todos}, nil
}
func (t *TodoServer) Read(ctx context.Context, req *models.GetPostRequest) (*models.GetResponse, error) {
	todo, ok := Todos[req.GetID()]
	if !ok {
		return nil, errors.New("not found")
	}
	return &models.GetResponse{Todo: todo}, nil
}
func (t *TodoServer) Create(ctx context.Context, req *models.ToDo) (*models.MutationResponse, error) {
	Todos[req.GetID()] = &models.ToDo{
		ID:   req.GetID(),
		Name: req.GetName(),
	}
	msg := req.GetID() + "successfully appended"
	return &models.MutationResponse{Success: msg}, nil
}
func (t *TodoServer) Update(ctx context.Context, req *models.UpdatePostRequest) (*models.MutationResponse, error) {
	Todos[req.GetID()] = &models.ToDo{
		ID:   req.GetID(),
		Name: req.GetName(),
	}
	msg := req.GetID() + "successfully appended"
	return &models.MutationResponse{Success: msg}, nil
}
func (t *TodoServer) Delete(ctx context.Context, req *models.DeletePostRequest) (*models.MutationResponse, error) {
	delete(Todos, req.GetID())
	msg := req.GetID() + "successfully deleted"
	return &models.MutationResponse{Success: msg}, nil
}
