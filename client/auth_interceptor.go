package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type AuthInterceptor struct {
	authClient AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor (
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
	) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient: *authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

//最重要的部分,添加拦截器以将令牌附加到请求上下文
//Unary
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func (
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--->unary interceptor: %s", method)

		//检查此方法是否需要身份验证,如果是这样就需要将访问令牌附加到上下文
		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		//如果该method不需要身份验证，那么就什么也不做
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
//Stream
func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--->stream interfector: %s", method)

		if interceptor.authMethods[method] {
			return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

//令牌刷新
func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	//此函数中应该确保第一次成功调用refreshToken()
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration		//记录令牌刷新钱还需要等多少时间
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

//无需调度即可刷新令牌
func (interceptor *AuthInterceptor) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)

	return nil
}