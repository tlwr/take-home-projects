test:
	go test -v $$(go list ./... | grep -v integration)
	go vet ./...

generate:
	go generate ./...

build:
	go build

docker:
	docker build -t truelayer-take-home-pokemon-api .

docker-run: docker
	docker run --rm -it -p 5000:5000
