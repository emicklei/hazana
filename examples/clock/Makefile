gen:
	protoc -I. --go_out=plugins=grpc:. clock.proto

server:
	go run server.go clock.pb.go

load:
	go run loadrun.go clock.pb.go -v -rps 2 -attack 20 