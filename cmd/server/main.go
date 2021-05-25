package main

import (
	"../../pb"
	"../../service"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

const (
	serverCertFile	 = 	"cert/server-cert.pem"
	serverKeyFile	 =	"cert/server-key.pem"
	clientCACertFile =	"cert/ca-cert.pem"
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
	//以下代码为开启双向TLS所需，创建新的证书池
	pemClientCA, err := ioutil.ReadFile(clientCACertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	//

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth: tls.RequireAndVerifyClientCert,	//这里开启双向TLS
		ClientCAs: certPool,			//开启双向TLS还需要提供手信任的CA证书列表
	}

	return credentials.NewTLS(config), nil
}

//定义运行gRPC server服务器的函数
func runGRPCServer(authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener,
) error {
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}
	if enableTLS {
		//获取TLS凭据
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS credentials: %w", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	grpcServer := grpc.NewServer(serverOptions...)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	//gRPC 反射, 然后可以通过evans测试grpc请求
	//调用reflection.Register(grpcServer)
	reflection.Register(grpcServer)

	log.Printf("start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	return  grpcServer.Serve(listener)
}


func runRESTServer(authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener,
	grpcEndpoint string,
) error {
	//创建一个新的HTTP请求多路复用器
	mux := runtime.NewServeMux()
	dialOption := []grpc.DialOption{grpc.WithInsecure()}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//开始从REST到gRPC的进程内转换
	//err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)	//auth-service, just support unary
	err := pb.RegisterLaptopServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, dialOption)	//auth0service, just support stream
	if err != nil {
		return err
	}
	err = pb.RegisterLaptopServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, dialOption)	//laptop-service
	if err != nil {
		return err
	}

	log.Printf("start REST server at %s, TLS = %t", listener.Addr().String(), enableTLS)

	//如果开启了TLS
	if enableTLS {
		return http.ServeTLS(listener, mux, serverCertFile, serverKeyFile)
	}
	//没有开启TLS
	return http.Serve(listener, mux)
}


func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSl/TLS")	//是否开启TLS
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")	//服务类型
	endPoint := flag.String("endpoint", "", "gRPC endpoint")
	flag.Parse()
	log.Printf("start server on port: %d, SSL/TLS = %t", *port, *enableTLS)

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

	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(authServer, laptopServer, jwtManager, *enableTLS, listener)
	} else {
		err = runRESTServer(authServer, laptopServer, jwtManager, *enableTLS, listener, *endPoint)
	}
	if err != nil {
		log.Fatal("cannot start server: ",err)
	}
}
