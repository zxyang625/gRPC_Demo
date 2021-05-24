# gRPC Demo

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
>    <u>当你运行之后会报错，说这种命令已经过时了，要用最新的命令运行</u>，像下面这样:
>
>    ```sh
>    $ protoc --go_out=. --go_opt=paths=source_relative \
>        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
>        helloworld/helloworld.proto
>    ```
>
>    但新命令完全不好使，而且网上能找到的教程大多都是使用旧命令，最新版不支持，<u>因此解决办法是不下载最新版的代码，在[https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)使用以前的版本</u>，我自己用的是v1.4.0。
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
> 2. 参照参考资料的配置方法，在Windows和Linux上运行时有个很大的坑。<u>由于Linux和Windows的文件路径分隔符是'/'和'\\'，因此在Windows上通过Git运行.sh会失败，解决方法是:将每个出现'/'命令的第一个‘/’改为"//"</u>。

