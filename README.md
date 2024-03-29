[![Build Status](https://github.com/stream1080/swinx/actions/workflows/go.yml/badge.svg)](https://github.com/stream1080/swinx/actions?query=branch%3Amaster) 
[![Go Report Card](https://goreportcard.com/badge/github.com/stream1080/swinx)](https://goreportcard.com/report/github.com/stream1080/swinx)
![license](https://img.shields.io/github/license/stream1080/swinx)
[![Go Reference](https://pkg.go.dev/badge/github.com/stream1080/swinx.svg)](https://pkg.go.dev/github.com/stream1080/swinx)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/stream1080/swinx)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/stream1080/swinx)

## 轻量级 TCP 服务器框架 swinx
Swinx 是一个 Go 语言实现的轻量级服务器框架，它可以快速开发基于网络通信的应用程序。

```       
                  _           
   ______      __(_)___  _  __
  / ___/ | /| / / / __ \| |/_/
 (__  )| |/ |/ / / / / />  <  
/____/ |__/|__/_/_/ /_/_/|_|  
                                
  ```
## 特性
- **轻量级**: Swinx 框架的代码量比较少，整个库只有几千行代码，易于学习和使用。
- **接口驱动**: Swinx 框架采用了面向接口编程的方式，提供了丰富的接口和组件，使你可以方便地实现网络通信相关的功能。
- **可插拔**: Swinx 框架的设计思路是把网络通信的各个组件分离开来，使你可以根据自己的需要选择合适的组件组合在一起。
- **可扩展**: Swinx 框架提供了丰富的扩展点，使你可以通过继承或实现接口来定制自己的功能。

## 面向接口
框架采用了面向接口编程的方式，提供了丰富的接口和组件，使你可以方便地实现网络
通信相关的功能。例如：
- 创建服务器
- 连接客户端
- 接收和发送数据

## 组件分离
框架的设计思路是把网络通信的各个组件分离开来，使你可以根据自己的需要选择合适
的组件组合在一起，例如：
- 定义网络协议;
- 定义消息编解码方式;
- 定义连接管理策略等。
