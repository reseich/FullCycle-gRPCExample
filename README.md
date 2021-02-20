sudo snap install protobuf --classic
go get google.golang.org/protobuf/cmd/protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go  
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

https://github.com/ktr0731/evans