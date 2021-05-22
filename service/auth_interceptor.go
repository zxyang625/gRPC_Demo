package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type AuthInterceptor struct {
	jwtManager *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

//Unary服务端拦截器
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func (
		ctx context.Context,
		req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	) (interface{}, error) {
	log.Println("-------> unary interceptor: ", info.FullMethod)

	err := interceptor.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
	}
}

//stream服务端拦截器
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func (
		srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
	) error {
	log.Println("-------> stream interceptor", info.FullMethod)

	err := interceptor.authorize(stream.Context(), info.FullMethod)
	if err != nil {
		return err
	}

	return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize (ctx context.Context, method string) error {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		return nil
	}

	//调用metadata包来获取请求的元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	//令牌为空，返回未经身份验证的token
	if len(values) == 0 {
		return status.Errorf(codes.Unimplemented, "authorization token is not provided")
	}

	//否则访问令牌应该储存在值的第一个元素中,调用jwtManager.Verify验证令牌并取回claims
	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		//找到了用户,返回nil
		if role == claims.Role {
			return nil
		}
	}

	//否则指出用户无权访问此RPC
	return status.Errorf(codes.PermissionDenied, "no permission to access this RPC")
}