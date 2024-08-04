.PHONY: server gen load

gen:
	protoc -I . \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	clock.proto
	
	cp clock.pb.go server
	cp clock_grpc.pb.go server

server:
	cd server && go run .

load:
	go run . -verbose -rps 2 -attack 20 