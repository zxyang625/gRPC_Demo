
gen:
	protoc --proto_path=./proto  --go_out=plugins=grpc:pb  proto/*.proto

clean:
	del pb\*.go

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 127.0.0.1:8080

#-cover衡量测试的代码覆盖率 -race检测代码中存在的竞争
test:
	go test -cover -race ./service

evans:
	evans -r -p 8080

.PHONY: gen client server clean test