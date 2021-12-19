binary:
	export GOOS=linux; go build -o target/hook   cmd/main.go
	export GOOS=linux; go build -o target/server examples/server/server.go

image: binary
	docker build . -t payall4u/nri-hook:latest

clean:
	rm -rf target