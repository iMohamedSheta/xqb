.PHONY: test
test: 
	go test -v ./...

.PHONY: lint
lint: 
	revive -formatter friendly ./...
