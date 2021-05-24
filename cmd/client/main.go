package main

import (
	"../../client"
	"../../pb"
	"../../sample"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
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
	username = "admin1"
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

//客户端只需要加载签署服务器的CA证书,因为客户需要验证真实性
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA' s certificate")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address","", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	//需要为身份验证客户端建立单独的连接，用于创建身份验证拦截器,将conn改为cc1
	cc1, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(tlsCredentials))
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
		grpc.WithTransportCredentials(tlsCredentials),
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
