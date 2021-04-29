.PHONY: vet lint test

vet:
	go vet ./...

lint:
	go get -u golang.org/x/lint/golint
	golint -set_exit_status ./...

run:
	docker-compose -f docker/docker-compose.yml up

run-%:
	go run ./go/cmd/$*/main.go

verify-pr: clean vet lint test build-all