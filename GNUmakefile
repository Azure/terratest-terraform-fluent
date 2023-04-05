TESTTIMEOUT=5m
TEST?=$$(go list ./... | grep -v 'vendor')

fmt:
	@echo "==> Fixing source code with gofmt..."
	find -name '*.go' | grep -v vendor | xargs gofmt -s -w

lint:
	golangci-lint run

test:
	go test $(TEST) $(TESTARGS) -run ^Test$(TESTFILTER) -timeout=$(TESTTIMEOUT)

# Create a test coverage report and launch a browser to view it
testcover:
	@echo "==> Generating testcover.out and launching browser..."
	if [ -f "coverage.txt" ]; then rm coverage.txt; fi
	go test $(TEST) $(TESTARGS) -timeout=$(TESTTIMEOUT) -coverprofile=coverage.txt -covermode=count
	go tool cover -html=coverage.txt

# Create a test coverage report in an html file
testcoverfile:
	@echo "==> Generating testcover html file..."
	if [ -f "coverage.txt" ]; then rm coverage.txt; fi
	if [ -f "coverage.html" ]; then rm coverage.html; fi
	go test $(TEST) $(TESTARGS) -timeout=$(TESTTIMEOUT) -coverprofile=coverage.txt -covermode=count
	go tool cover -html=coverage.txt -o=coverage.html

tools:
	go install mvdan.cc/gofumpt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.46.2

.PHONY: fmt lint test testcover testcoverfile tools
