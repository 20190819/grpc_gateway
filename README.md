### grpc gateway demo

- 安装
```code
go get github.com/golang/protobuf/proto

go get google.golang.org/grpc

go get github.com/golang/protobuf/protoc-gen-go

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

- 编译命令
```code
protoc -I. -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  --go_out=plugins=grpc:. ./proto/helloworld.proto
```
```code
protoc -I. -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis   --grpc-gateway_out=logtostderr=true:. ./proto/helloworld.proto
```

- 模拟请求
```code
curl -v -X POST "http://localhost:8082/v1/helloworld/getfile" -d '{"value":"server.go"}'

curl -v -X POST "http://localhost:8082/v1/helloworld/getfile" -d '{"value":"server"}'
```