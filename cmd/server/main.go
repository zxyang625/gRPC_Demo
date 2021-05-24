package main

import (
	"../../pb"
	"../../service"
	"crypto/tls"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)


func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username string, password string, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}


//简单起见，secretKey和tokenDuration设置为常量，实际应该从环境变量或者配置文件中加载
const (
	secretKey = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/pb.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop" : {"admin"},
		laptopServicePath + "uploadImage" : {"admin"},
		laptopServicePath + "RateLaptop" : {"admin", "user"},
	}
}

//加载TLS凭据,对于服务端需要加载服务器的证书和私钥,因此采用tls.LoadX509KeyPair函数
//从cert文件夹加载server-cert.pem和server-key.pem文件
//然后使用服务器证书制作一个tls配置对象
//最后使用该凭据并返回给呼叫者
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth: tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port: %d", *port)

	userStore := service.NewInMemoryUserStore()
	//测试新建用户
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	//获取TLS凭据
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
		)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	//gRPC 反射, 然后可以通过evans测试grpc请求
	//调用reflection.Register(grpcServer)
	reflection.Register(grpcServer)

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
