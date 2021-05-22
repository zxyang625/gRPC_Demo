# gRPC Demo

<h2>记录gRPC中趟过的坑(持续更新)</h2>


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

