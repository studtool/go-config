TEST_DIR := $(PWD)/tests
TOOLS_DIR := $(PWD)/tools

COVERAGE_PROFILE_FILE := $(TEST_DIR)/c.out
COVERAGE_REPORT_FILE := $(TEST_DIR)/c.html

TEST_FLAGS := -mod=vendor -v -race -coverprofile='$(COVERAGE_PROFILE_FILE)'
COVER_FLAGS := -html='$(COVERAGE_PROFILE_FILE)' -o '$(COVERAGE_REPORT_FILE)'

fmt:
	go fmt ./...

dep:
	go get -u && go mod tidy && go mod vendor && go mod verify

test: mkdirs bootstrap dep
	go test $(TEST_FLAGS) ./...
	go tool cover $(COVER_FLAGS)

clean:
	rm $(COVERAGE_PROFILE_FILE)
	rm $(COVERAGE_REPORT_FILE)

lint: mkdirs bootstrap
	$(TOOLS_DIR)/golangci-lint run

mkdirs:
	mkdir -p $(TEST_DIR)
	mkdir -p $(TOOLS_DIR)

bootstrap:
	./install_tools.sh '$(TOOLS_DIR)'
