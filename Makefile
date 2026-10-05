GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),arm64)
	DARWIN_GOARCH=arm64
else ifeq ($(UNAME_M),x86_64)
	DARWIN_GOARCH=amd64
else
	DARWIN_GOARCH=$(shell go env GOARCH)
endif

GOARCH ?= $(shell go env GOARCH)
GOOS=$(shell go env GOOS)
LINUX_GOARCH ?= amd64

BASE_PATH := $(shell pwd)
BUILD_PATH = $(BASE_PATH)/build
WEB_PATH=$(BASE_PATH)/frontend
ASSERT_PATH= $(BASE_PATH)/core/cmd/server/web/assets

CORE_PATH=$(BASE_PATH)/core
CORE_MAIN=$(CORE_PATH)/cmd/server/main.go
CORE_NAME=1panel-core

AGENT_PATH=$(BASE_PATH)/agent
AGENT_MAIN=$(AGENT_PATH)/cmd/server/main.go
AGENT_NAME=1panel-agent

MACPANEL_PATH=$(BASE_PATH)/cmd/macpanel
MACPANEL_NAME=macpanel

# Embed build order (required for correct frontend assets in core/macpanel):
#   1. make clean_assets          (optional)
#   2. make build_frontend        Vue build -> core/cmd/server/web/assets/
#   3. make build_for_darwin      go embed via core dependency in macpanel
#   4. ./scripts/mac/start.sh     run unified binary

clean_assets:
	rm -rf $(ASSERT_PATH)

upx_bin:
	upx $(BUILD_PATH)/$(CORE_NAME)
	upx $(BUILD_PATH)/$(AGENT_NAME)

build_frontend:
	cd $(WEB_PATH) && npm install && npm run build:pro

build_core_on_linux:
	cd $(CORE_PATH) \
	&& CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)

build_agent_on_linux:
	cd $(AGENT_PATH) \
    && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

# Native macOS unified binary (core + agent in one process)
build_for_darwin: build_frontend
	mkdir -p $(BUILD_PATH)
	cd $(MACPANEL_PATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=$(DARWIN_GOARCH) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(MACPANEL_NAME) .

# Cross-compile Linux server binaries (separate core + agent)
build_for_linux: build_frontend
	mkdir -p $(BUILD_PATH)
	cd $(CORE_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=$(LINUX_GOARCH) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)
	cd $(AGENT_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=$(LINUX_GOARCH) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

build_macpanel: build_for_darwin

# Legacy cross-compile to Linux (kept for compatibility)
build_core_on_darwin:
	cd $(CORE_PATH) \
	&&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w'  -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)

build_agent_on_darwin:
	cd $(AGENT_PATH) \
    &&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w'  -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

build_all: build_frontend build_core_on_linux build_agent_on_linux

build_on_local: clean_assets build_macpanel

.PHONY: clean_assets upx_bin build_frontend build_core_on_linux build_agent_on_linux \
	build_for_darwin build_for_linux build_macpanel build_core_on_darwin \
	build_agent_on_darwin build_all build_on_local
