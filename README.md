# gRPC Demo

<h2>项目说明</h2>

这是学习gRPC而实现的一个Demo，这个Demo的主要功能就是通过gRPC实现对笔记本电脑参数的创建，储存，评分等功能。期间陆陆续续完善例如实现四种gRPC调用，通过evans测试gRPC请求，通过JWT实现拦截器和身份验证，通过Openssl生成自签证书完成双向TLS验证，采用nginx实现负载均衡以及配置gRPC-Gateway来实现REST到gRPC的转换。

├─.idea  
├─cert						//存放证书以及密钥	    	  
├─client					 //用户登录以及验证令牌  
├─cmd  
│  ├─client				 //客户端发送请求  
│  └─server			   //服务器响应请求  
├─pb						  //编译的.pb.go文件  
├─proto					 //proto文件  
│  └─google			   //谷歌http.api批注  
│      ├─api  
│      └─rpc  
├─sample				   //实例生成  
├─serializer				//序列化和反序列化  
├─service					//gRPC服务  
├─swagger				 //swagger反射json文件  
└─tmp						 //临时文件    



<h2>细数gRPC中的坑(持续更新)</h2>

- <h4>protobuf这个包到底在哪?github.com/protobuf? github.com/protocolbuffers/protobuf?  google.golang.org/protobuf?</h4>

> 简单来说，[github.com/protocolbuffers/protobuf](https://github.com/protocolbuffers/protobuf)是为各种语言配置protobuf依赖的说明大纲。[github.com/golang/protobuf](https://github.com/golang/protobuf)是最开始go对protobuf的支持的项目，然后这个项目被google接管，因此项目地址也变成了[google.golang.org/protobuf](google.golang.org/protobuf)，不过这个google的项目又把代码托管到了github上，因此最终项目的地址又变成了[github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)。这就是为什么protobuf这个包在三个地方都有。  

> 网上对于go如何配置protobuf和gRPC都有相关的文档资料，不过其中大部分都只告诉了你要怎么做，由于版本改变的原因，经常出现一些配置不成功或者命令不对的坑，我这里就来说明一下相关的坑:
>
> 1. 关于protoc.exe以及proto-gen-go.exe
>
>    protoc 命令来自于 [https://github.com/google/protobuf](https://github.com/google/protobuf)，可以产生序列化和反序列化的代码，无go相关代码。protoc-gen-go插件则来自于[github.com/golang/protobuf/protoc-gen-go](https://github.com/golang/protobuf/protoc-gen-go)， 可以产生go相关代码， 除上述序列化和反序列化代码之外， 还增加了一些通信公共库。
>
>    当然，这两个库也可以在新版的[https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)上安装，不过坑也因此产生。 如果你从新地址上拷代码下来，再输入protoc命令，例如:
>
>    ```
>    $ protoc --go_out=./go1/ ./proto/my.proto
>    $ protoc --go_out=plugins=grpc:./go2/ ./proto/my.proto
>    ```
>
>    **当你运行之后会报错，说这种命令已经过时了，要用最新的命令运行**，像下面这样:
>
>    ```sh
>    $ protoc --go_out=. --go_opt=paths=source_relative \
>        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
>        helloworld/helloworld.proto
>    ```
>
>    但新命令完全不好使，而且网上能找到的教程大多都是使用旧命令，最新版不支持，**因此解决办法是不下载最新版的代码，在[https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)使用以前的版本**，我自己用的是v1.4.0。
>
> 2. 关于gRPC库依赖下载
>
>    安装官方安装命令：
>    `go get google.golang.org/grpc`
>    是安装不起的，会报：
>
>    `package google.golang.org/grpc: unrecognized import path "google.golang.org/grpc"(https fetch: Get https://google.golang.org/grpc?go-get=1: dial tcp 216.239.37.1:443: i/o timeout)`
>
>    原因google.golang.org被墙了，并且这个代码已经转移到github上面了，但是代码里面的包依赖还是没有修改，还是google.golang.org。也因此，仅仅从github上面下载代码还是不够。需要将对应github.com的名字改为google.golang.org。
>
>    因此解决办法是:
>
>    ```.sh
>    git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc  
>    git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net  
>    git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text  
>    go get -u github.com/golang/protobuf/proto
>    go get -u github.com/golang/protobuf/proto/protoc-gen-go} 
>    git clone https://github.com/google/go-genproto.git $GOPATH/src/google.golang.org/genproto  
>                               
>    cd $GOPATH/src/  
>    go install google.golang.org/grpc 
>    ```
>
> 

- <h4>evans该怎么用?</h4>


> [evans](https://github.com/ktr0731/evans)是一款带命令行补全的交互式客户端，用于向运行的服务器发送请求，以此免去测试文件，提高开发效率的工具。
> 对于windows来说，需要下载release中对应系统的压缩包，最后将解压出来的.exe文件放入GOPATH中。然后在命令行输入参数即可测试自己的gRPC server。
> 不过其中可能会遇到像我一样运行没用的情况，在翻了一天的谷歌无果后我仔细地看了evans中的方方面面。最后发现是由于版本原因和命令改变导致运行evans无法成功。因此需要注意的点如下:
>
> 1. 确保下载的evans是对应当前系统架构的release
> 2. 确保server已经启动
> 3. 确保已在server中开启了reflection反射服务
>
> 然后就可以在cmd中输入一下命令开启evans了(假设端口为8080,默认端口为50051):
>
> `$ evans -r -p 8080`

- <h4>关于使用openssl生成自签证书</h4>


> 关于openssl命令以及证书生成的[参考资料](https://blog.csdn.net/qq_30145355/article/details/113279539)。
>
> ca-key.pem：CA的私钥  
> ca-cert.pem：CA的证书文件  
> ca-cert.srl：保存证书的唯一序列号  
> client-ext.cnf：客户端证书扩展信息。  
> client-req.pem：客户端的证书签名请求（CSR）文件  
> client-key.pem：客户端的私钥  
> client-cert.pem：客户端的证书  
> server-ext.cnf：服务器证书扩展信息  
> server-req.pem：服务器的证书签名请求（CSR）文件  
> server-key.pem：服务器的私钥  
> server-cert.pem：服务器的证书  

> 关于openssl比较坑的地方:
>
> 1. openssl的二进制文件在原网址上移除了，而如果要自己配置会相当麻烦
>
> 2. 参照参考资料的配置方法，在Windows和Linux上运行时有个很大的坑。**由于Linux和Windows的文件路径分隔符是'/'和'\\'，因此在Windows上通过Git运行.sh会失败，解决方法是:将每个出现'/'命令的第一个‘/’改为"//"**。
>
> 3. **在localhost中使用自签证书时，server-ext.cnf和client-ext.cnf中必须要指定IP地址作为证书使用者备用名称(SAN)的扩展**，例如:
>
>    `subjectAltName=DNS:*.study.com,DNS:*.study.org,IP:127.0.0.1`
>
>    这是因为如果IP不对应server地址或者缺少IP，会导致TLS无法验证证书而握手失败。而在实际生产中是可以忽略IP的，因为可以改用域名。

- <h4>在gRPC中使用nginx</h4>

> nginx支持HTTP2.0，因此可以在gRPC中使用。[nginx](http://nginx.org/en/download.html)中各版本的区别:
>
> Mainline version：Nginx 目前主力在做的版本，可以认为是主线版本、开发版  
> Stable version：最新的稳定版，生产环境上建议使用这个版本  
> Legacy versions：稳定版的历史版本集合  
>
> **需要注意的是:nginx的版本号并不是数字越大版本越高**，例如[nginx-1.8.1]已经是相当老旧的版本了，使用老版本是跑步起来的，我自己用的是1.20.0。
>
> **nginx对HTTP2的支持需要在conf/nginx.conf修改配置,例如我在conf文件夹下新建了cert文件夹用于存放.pem文件**，以下是未开启TLS时的配置:
>
> ```
>  worker_processes  1;
>   error_log  logs/error.log;
> 
>   events {
>       worker_connections  10;
>   }
> 
>   http {
>   access_log  logs/access.log;
> 
>   upstream pcbook_services {
>   server 127.0.0.1:50051;
>   server 127.0.0.1:50052;
>   }
> 
>   server {
>       listen       8080 ssl http2;
>       #告诉nginx证书和密钥的位置
>       ssl_certificate cert/server-cert.pem;
>       ssl_certificate_key cert/server-key.pem;
>       ssl_client_certificate cert/ca-cert.pem;
>   ssl_verify_client on;   #开启告诉nginx验证客户端发送证书的真实性
> 
>   location / {
>  	grpc_pass grpc://pcbook_services;
>   }
> ```

> 1. 在Windows10 下使用Nginx可能会出现问题:
>
>`nginx: [emerg] BIO_new_file("./conf/cert/nginx.pem") failed (SSL: error:02001003:system library:fopen:No such process:fopen(’./conf/cert/nginx.pem’,‘r’) error:2006D080:BIO routines:BIO_new_file:no such file)`
> 
>这个问题的出现代表nginx配置文件中配置了ssl协议，但nginx确没有相应的证书文件。所以应当检查证书文件的路径是否正确，也可以考虑加入 -c 显式指定配置的conf路径:
> 
>```
> start nginx -c ./conf/nginx.exe
> ```

> 2. 此外，**如果在服务器和客户端之间启用了TLS连接，那么在加入nginx之后启动客户端发送请求肯定是会失败的。因为尽管nginx和客户端的TLS握手成功，但是nginx和后端服务器的TLS握手仍会失败，这是由于服务端想启用TLS连接，但是nginx向服务器的发起的连接使用的仍然是不安全的连接。**解决办法如下:
>    1. 只启用服务端TLS
>
>       1) 在nginx.conf中将grpc方案修改为grpcs:
>
>    ```
>    location / {
>     	grpc_pass grpcs://pcbook_services;	#启用TLS需要将grpc改为grpcs
>      }
>    ```
>
>    ​		2) 修改server.go中代码，设置config的ClientAuth字段为 `tls.NoClientCert`
>
>    **注意:这个方法在只使用服务端TLS的时候才有效。**
>
>    2. 启用双向TLS
>
>       1) 在nginx.conf中将grpc方案修改为grpcs:
>
>    ```
>    location / {
>     	grpc_pass grpcs://pcbook_services;	#启用TLS需要将grpc改为grpcs
>      }
>    ```
>
>    ​	   2) 在nginx.conf中指定服务器证书的位置:
>
>    ```
>     location / {
>    	grpc_pass grpcs://pcbook_services; 	#启用TLS需要将grpc改为grpcs
>          
>    	#双向TLS需要指定nginx发送给上游服务器的证书的位置
>    	grpc_ssl_certificate cert/server-cert.pem;
>    	grpc_ssl_certificate_key cert/server-key.pem;
>     }
>    ```
>
>    **提示:实际上应该也为nginx生成另一对证书和密钥，以满足nginx和服务器之间的TLS，这里我简单使用已有的证书和密钥。**

- <h4>关于grpc-gateway的使用</h4>

> [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)是protoc的一个插件。它能读取gRPC服务定义，并生成一个反向代理服务器，将RESTFUL JSON API转换为gRPC。举例来说，假定我们有一个项目需求，希望**用RPC作为内部API的通讯**，同时也**想对外提供Restful API**，写两套又太繁琐不符合，于是可以利用gRPC以及gRPC Gateway。

> gRPC-Gateway的坑点在于其配置过程，很多资料由于年代久远，未能及时更新gRPC-Gateway的版本迭代，再加上官方文档写的属实过于偷懒，因此配置过程中出现错误也很正常。
>
> **首先需要说明的是，最新版gRPC-Gateway v2不再支持GOPATH，只支持GO MOD!**
>
> 所以如果采用GOPATH管理package的话必然会遇到如下错误:
>
> ```
> cannot find package "github.com/grpc-ecosystem/grpc-gateway/v2/utilities" in any of:
>  E:\Golang\Go\src\github.com\grpc-ecosystem\grpc-gateway\v2\utilities (from $GOROOT)
>  C:\Users\Lenovo\go\src\github.com\grpc-ecosystem\grpc-gateway\v2\utilities (from $GOPATH)
> 
> ```
>
> 解决办法是要么就采用GO MOD管理package，要么继续使用v1的最新版[v1.15.1](https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.15.1)，不过v1版本已经停止更新了,**在使用v1版本时要注意其他包和工具的版本，不然出现以下警告是必然的**:
>
> ```
> WARNING: Package "github.com/golang/protobuf/protoc-gen-go/generator" is deprecated.
>      A future release of golang/protobuf will delete this package,
>      which has long been excluded from the compatibility promise.
> 
> ```
>
> 排除这个坑之后，按照官方文档的方法来就没有问题了。不过有些命令也已经改变了，详情可以参考[gRPC-Gateway v2](https://grpc-ecosystem.github.io/grpc-gateway/docs/development/grpc-gateway_v2_migration_guide/)关于新版和旧版的区别。

> 如何实现从REST到gRPC的转换?
>
> 通过**RegisterXXXServiceHandleServer()** 函数实现从REST到gRPC的进程内转换，这样就不需要单独的gRPC服务器通过网络调用来服务REST发来的请求。不过此函数只支持Unary RPC，对于stream RPC，则需要使用**RegisterXXXServiceHandleFromEndpoint()** 函数，这个函数将传入的RESTful请求转换为gRPC格式，并在指定的端口上调用相应的RPC。

