package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RegisterReflection registers gRPC reflection with the server
func RegisterReflection(server *grpc.Server) {
	reflection.Register(server)
}

// RegisterReflectionWithServices registers gRPC reflection with specific services
func RegisterReflectionWithServices(server *grpc.Server, services ...string) {
	reflection.Register(server)
}
