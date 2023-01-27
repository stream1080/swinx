package pack

import (
	"sync"

	"github.com/stream1080/swinx/face"
)

type pack_factory struct{}

var (
	pack_once       sync.Once
	factoryInstance *pack_factory
)

// 创建一个拆包封包的实例，单例
func Factory() *pack_factory {
	pack_once.Do(func() {
		factoryInstance = new(pack_factory)
	})

	return factoryInstance
}

func (f *pack_factory) NewPack(king string) face.IDataPack {
	var dataPack face.IDataPack

	switch king {
	case face.SwinxDataPack:
		dataPack = NewDataPack()
	case face.CustomDataPack:
		// 自定义封包拆包方式
	default:
		dataPack = NewDataPack()
	}

	return dataPack
}
