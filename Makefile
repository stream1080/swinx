SHELL = /bin/bash

help:
	@echo "make 运行程序"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试"
	@echo "make server 启动服务端"
	@echo "make client1 启动客户端1"
	@echo "make client1 启动客户端2"

build: clean
	@go build

clean:
	@go clean

test:
	@go test

server:
	@go run main.go -t server

client1:
	@go run main.go -t client1

client2:
	@go run main.go -t client2