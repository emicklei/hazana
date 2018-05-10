.PHONY: server gen load

gen:
	protoc -I. --go_out=plugins=grpc:. clock.proto
	cp clock.pb.go server

server:
	cd server && go run server.go clock.pb.go

load:
	go run loadrun.go clock.pb.go -verbose -rps 2 -attack 20 