.PHONY: test
test: 
	go test -v --cover ./...

.PHONY: lint
lint: 
	revive -formatter friendly ./...
