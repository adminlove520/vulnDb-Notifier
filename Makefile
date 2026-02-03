build:
	go build -o vulnDb-Notifier cmd/CVENotifier/main.go

run:
	go run cmd/CVENotifier/main.go