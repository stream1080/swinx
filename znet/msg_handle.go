package znet

import (
	"fmt"
	"strconv"

	"github.com/stream1080/swinx/conf"
	"github.com/stream1080/swinx/face"
)

type MsgHandle struct {
	HandleMap      map[uint32]face.IRouter // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                  // 业务工作 Worker 池的数量
	TaskQueue      []chan face.IRequest    // Worker 负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	poolSize := conf.ServerConfig.WorkerPoolSize
	return &MsgHandle{
		HandleMap:      map[uint32]face.IRouter{},
		WorkerPoolSize: poolSize,
		TaskQueue:      make([]chan face.IRequest, poolSize),
	}
}

// 以非阻塞方式处理消息
func (h *MsgHandle) DoMsgHandler(request face.IRequest) {

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
func (h *MsgHandle) AddRouter(msgId uint32, router face.IRouter) {

	// 判断当前 msg 绑定的 router 处理方法是否已经存在
	_, ok := h.HandleMap[msgId]
	if ok {
		panic("repeated router , msgId:" + strconv.Itoa(int(msgId)))
	}

	// 添加 msgId 与 router 的绑定关系
	h.HandleMap[msgId] = router
	fmt.Println("add router msgId:", msgId)
}

// 启动 worker 工作池
func (h *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		// 给当前worker对应的任务队列开辟空间
		h.TaskQueue[i] = make(chan face.IRequest, conf.ServerConfig.MaxWorkerTaskSize)
		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go h.StartWorker(i, h.TaskQueue[i])
	}
}

func (h *MsgHandle) StartWorker(workerId int, taskQueue chan face.IRequest) {
	fmt.Println("WorkerId:", workerId, " is started")
	// 等待队列中的消息
	for req := range taskQueue {
		// 有消息则取出队列的 Request，并执行绑定的业务方法
		h.DoMsgHandler(req)
	}
}

// 将消息交给 TaskQueue ,由 worker 进行处理
func (h *MsgHandle) SendMsgTaskQueue(request face.IRequest) {
	// 根据 ConnId 来分配当前的连接应该由哪个 worker 负责处理
	// 轮询的平均分配法则
	// 得到需要处理此条连接的 workerId
	workerId := request.GetConn().GetConnId() % h.WorkerPoolSize
	fmt.Printf("add connId: %d, request msgId: %d, to workerId: %d \n", request.GetConn().GetConnId(), request.GetMsgId(), workerId)

	// 将请求消息发送给任务队列
	h.TaskQueue[workerId] <- request
}
