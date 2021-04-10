SHELL := /bin/bash

# ==============================================
run:
	go run cmd/dashboard-api/main.go

run-admin:
	go run cmd/dashboard-admin/main.go

test:
	go test -v ./... -count=1
	staticcheck ./...
	
tidy:
	go mod tidy
	go mod vendor