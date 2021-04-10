SHELL := /bin/bash

# ==============================================
run:
	go run app/dashboard-api/main.go

run-admin:
	go run app/dashboard-admin/main.go

test:
	go test -v ./... -count=1
	staticcheck ./...
	
tidy:
	go mod tidy
	go mod vendor