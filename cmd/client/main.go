package main

import (
	"../../client"
	"../../pb"
	"../../sample"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)?")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}
		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz: 2.5,
		MinRam: &pb.Memory{
			Value: 8,
			Uint: pb.Memory_GIGABYTE,
		},
	}
	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/laptop.jpg")
}

const (
	username = "user1"
	password = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/pb.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop" : true,
		laptopServicePath + "uploadImage" : true,
		laptopServicePath + "RateLaptop" : true,
	}
}

func main() {
	serverAddress := flag.String("address","", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	//需要为身份验证客户端建立单独的连接，用于创建身份验证拦截器,将conn改为cc1
	cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(":cannot dial server:", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	//用客户端创建一个新的拦截器
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	//如果没有就Dial服务器来创建另一个连接
	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := client.NewLaptopClient(cc2)
	testRateLaptop(laptopClient)
	//for i := 0; i < 10; i++ {
	//	createLaptop(laptopClient, sample.NewLaptop())
	//}
	//filter := &pb.Filter{
	//	MaxPriceUsd: 3000,
	//	MinCpuCores: 4,
	//	MinCpuGhz: 2.5,
	//	MinRam: &pb.Memory{
	//		Value: 8,
	//		Uint: pb.Memory_GIGABYTE,
	//	},
	//}
	//searchLaptop(laptopClient, filter)

}
