package znet

import (
	"fmt"
	"strconv"

	"github.com/stream1080/swinx/face"
)

type MsgHandle struct {
	HandleMap map[uint32]face.Router // 存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		HandleMap: map[uint32]face.Router{},
	}
}

// 以非阻塞方式处理消息
func (h *MsgHandle) DoMsgHandler(request face.Request) {

	// 判断该 request 的 router 是否存在
	handler, ok := h.HandleMap[request.GetMsgId()]
	if !ok {
		fmt.Println("router msgId:", request.GetMsgId(), " is not found!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体处理逻辑
func (h *MsgHandle) AddRouter(msgId uint32, router face.Router) {

	// 判断当前 msg 绑定的 router 处理方法是否已经存在
	_, ok := h.HandleMap[msgId]
	if ok {
		panic("repeated router , msgId:" + strconv.Itoa(int(msgId)))
	}

	// 添加 msgId 与 router 的绑定关系
	h.HandleMap[msgId] = router
	fmt.Println("add router msgId:", msgId)
}
