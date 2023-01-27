package face

// 直接面向 TCP 连接中的数据流, 为传输数据添加头部信息，用于处理 TCP 粘包问题
type IDataPack interface {
	GetHeadLen() uint32                // 获取包头长度
	Pack(msg IMessage) ([]byte, error) // 封包方法
	UnPack([]byte) (IMessage, error)   // 拆包方法
}

const (
	SwinxDataPack  string = "swinx_pack"  // Swinx 标准封包和拆包方式
	CustomDataPack string = "custom_pack" // 自定义封包方式

	SwinxMessage  string = "swinx_message"  // Swinx 默认标准报文协议格式
	CustomMessage string = "custom_message" // 自定义报文协议格式
)
