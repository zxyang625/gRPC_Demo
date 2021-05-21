package main

import (
	"../../pb"
	"../../service"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

//Unary服务端拦截器
func unaryInterceptor (
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	) (interface{}, error) {
	log.Println("-------> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

//stream服务端拦截器
func streamInterceptor (
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
	) error {
	log.Println("-------> stream interceptor", info.FullMethod)
	return handler(srv, stream)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port: %d", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
		)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	//gRPC 反射, 然后可以通过evans测试grpc请求
	//调用reflection.Register(grpcServer)
	//reflection.Register(grpcServer)

	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("cannot start server", err)
	}

}
