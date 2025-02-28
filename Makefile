build:
	go build -o bin/chat main.go

run-server:
	bin/chat

run-client:
	bin/chat client