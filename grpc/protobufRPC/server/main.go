package main

import (
	"github.com/678g0/golangbook/grpc/protobufRPC/service"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 通过接口约束HelloService服务
var _ service.HelloService = (*HelloService)(nil)

type HelloService struct{}

func (p *HelloService) Hello(req *service.Request, resp *service.Reply) error {
	resp.Value = "hello:" + req.Value
	return nil
}

func main() {
	// 把我们的对象注册成一个rpc的 receiver
	// 其中rpc.Register函数调用会将对象类型中所有满足RPC规则的对象方法注册为RPC函数，
	// 所有注册的方法会放在“HelloService”服务空间之下
	_ = rpc.RegisterName("HelloService", new(HelloService))

	// 然后我们建立一个唯一的TCP链接，
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务。
	// 没Accept一个请求，就创建一个goroutie进行处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// // 前面都是tcp的知识, 到这个RPC就接管了
		// // 因此 你可以认为 rpc 帮我们封装消息到函数调用的这个逻辑,
		// // 提升了工作效率, 逻辑比较简洁，可以看看他代码
		// go rpc.ServeConn(conn)

		// 代码中最大的变化是用rpc.ServeCodec函数替代了rpc.ServeConn函数，
		// 传入的参数是针对服务端的json编解码器
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}