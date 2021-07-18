package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"yangliang4488/grpc_gateway/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	LISTEN_GRPC = ":9092"
	LISTEN_HTTP = ":8082"
)

var err error

type MyServer struct{}

func (my *MyServer) Echo(ctx context.Context, req *proto.HRequest) (rep *proto.HResponse, err error) {

	return &proto.HResponse{
		Data: "hello: " + req.Value,
	}, nil
}

func (my *MyServer) DownloadFile(req *proto.HRequest, rep proto.YourService_DownloadFileServer) error {
	if req.Value == "" {
		return errors.New("query params empty")
	}
	fileBuffer, err := ioutil.ReadFile(req.Value)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 每次发送字节
	byteOnce := 50
	// 第几次发送
	index := 0
	// 共计发送字节
	cnt := 0

	// 循环发送
	for {
		if index*byteOnce+byteOnce >= len(fileBuffer) {
			rep.Send(&proto.FileBinary{Data: fileBuffer[index*byteOnce:]})
			cnt += len(fileBuffer) - index*byteOnce
			break
		} else {
			rep.Send(&proto.FileBinary{Data: fileBuffer[index*byteOnce : index*byteOnce+byteOnce]})
			cnt += byteOnce
		}
		index++
	}
	fmt.Printf("cnt=%d\n", cnt)
	return nil
}

func main() {
	// grpc 服务
	lis, err := net.Listen("tcp", LISTEN_GRPC)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	proto.RegisterYourServiceServer(grpcServer, &MyServer{})
	go grpcServer.Serve(lis)
	fmt.Println("start grpc server", LISTEN_GRPC)

	// http 服务
	mux := runtime.NewServeMux()
	err = proto.RegisterYourServiceHandlerFromEndpoint(context.Background(), mux, LISTEN_GRPC, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("start http server", LISTEN_HTTP)
	http.ListenAndServe(LISTEN_HTTP, mux)
}
