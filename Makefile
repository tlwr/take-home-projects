.PHONY: integration

test:
	go test -v $$(go list ./... | grep -v integration)

integration:
	go test -v ./integration
