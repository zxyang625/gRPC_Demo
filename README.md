# gRPC Demo

<h2>细数gRPC中的坑(持续更新)</h2>


- <h4>protobuf这个包到底在哪?github.com/protobuf or google/protobuf or google.golang.org/protobuf?</h4>

- <h4>evans该怎么用?</h4>
  
  > [evans](https://github.com/ktr0731/evans)是一款带命令行补全的交互式客户端，用于向运行的服务器发送请求，以此免去测试文件，提高开发效率的工具。
  >
  > 对于windows来说，需要下载release中对应系统的压缩包，最后将解压出来的.exe文件放入GOPATH中。然后在命令行输入参数即可测试自己的gRPC server。
  >
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
> ca-cert.pem：CA的证书文件。
> ca-cert.srl：保存证书的唯一序列号。
> client-ext.cnf：客户端证书扩展信息。
> client-req.pem：客户端的证书签名请求（CSR）文件。
> client-key.pem：客户端的私钥。
> client-cert.pem：客户端的证书。
> server-ext.cnf：服务器证书扩展信息。
> server-req.pem：服务器的证书签名请求（CSR）文件。
> server-key.pem：服务器的私钥。
> server-cert.pem：服务器的证书。
>
> 关于openssl比较坑的地方:
>
> 1. openssl的二进制文件在原网址上移除了，而如果要自己配置会相当麻烦
> 2. 参照参考资料的配置方法，在Windows和Linux上运行时有个很大的坑。<u>由于Linux和Windows的文件路径分隔符是'/'和'\\'，因此在Windows上通过Git运行.sh会失败，解决方法是:将每个出现'/'命令的第一个‘/’改为"//"</u>。

