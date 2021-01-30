test:
	go test -v $$(go list ./... | grep -v integration)
	go vet ./...

generate:
	go generate ./...

build:
	go build
