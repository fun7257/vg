BINARY_NAME=vg
BUILD_DIR=.tmp/build

.PHONY: all build clean run install-hooks

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

clean:
	rm -rf $(BUILD_DIR)

install-hooks:
	@./scripts/install-hooks.sh
