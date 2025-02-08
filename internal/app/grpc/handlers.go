package grpc

import (
	"context"
	"errors"
	proto "github.com/Tairascii/google-docs-protos/gen/go/user"
	"github.com/Tairascii/google-docs-user/internal/app"
	"github.com/Tairascii/google-docs-user/internal/app/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Grpc struct {
	proto.UnimplementedUserServer
	DI *app.DI
}

func NewGrpc(DI *app.DI) *Grpc {
	return &Grpc{DI: DI}
}

type GrpcServer struct {
	server *grpc.Server
	lis    net.Listener
}

func NewGrpcServer(port string, DI *app.DI) (*GrpcServer, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	grpcStruct := NewGrpc(DI)
	proto.RegisterUserServer(grpcServer, grpcStruct)
	return &GrpcServer{server: grpcServer, lis: lis}, nil
}

func (s *GrpcServer) Start() error {
	log.Printf("gRPC server listening on %s", s.lis.Addr())
	return s.server.Serve(s.lis)
}

func (g *Grpc) IdByEmail(ctx context.Context, in *proto.IdByEmailRequest) (*proto.IdByEmailResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}

	id, err := g.DI.UseCase.User.IdByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.IdByEmailResponse{Id: id.String()}, nil
}
