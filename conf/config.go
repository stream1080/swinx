package conf

import (
	"encoding/json"
	"os"
)

type Config struct {
	//TcpServer face.Server // Tcp 服务器
	Name string // 名称
	Host string // 主机地址
	Port int    // 端口

	Version        string // 服务的版本
	MaxConn        int    // 最大连接数
	MaxPackageSize uint32 // 一次请求的最大数据包
}

// 全局对象
var ServerConfig *Config

// 加载配置
func (s *Config) LoadConfig() {
	conf, err := os.ReadFile("conf/conf.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(conf, &ServerConfig)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 定义默认配置
	ServerConfig = &Config{
		Name:           "Tcp-Server",
		Host:           "0.0.0,0",
		Port:           8888,
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 使用自定义配置
	ServerConfig.LoadConfig()
}
