SHELL = /bin/bash

help:
	@echo "make 运行程序件"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试"

build: clean
	go build

clean:
	go clean

test:
	go test

server:
	cd examples && go run server.go

client1:
	cd examples && go run client1.go

client2:
	cd examples && go run client2.go