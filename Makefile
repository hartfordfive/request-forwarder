GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=request-forwarder
GO_DEP_FETCH=govendor fetch 
UNAME=$(shell uname)
BUILD_DIR=build/
GITHASH=$(git rev-parse --verify HEAD)
BUILDDATE=$(date +%Y-%m-%d)
VERSION=$(shell sh -c 'cat VERSION.txt')
PACKAGE_BASE="github.com/hartfordfive/prometheus-ldap-sd"

ifeq ($(UNAME), Linux)
	OS=linux
endif
ifeq ($(UNAME), Darwin)
	OS=darwin
endif
ARCH=amd64

ifeq ($(ADD_VERSION_OS_ARCH), 1)
	BINARY_NAME=$(BINARY_NAME)-$(VERSION)-$(OS)-$(ARCH)
endif

all: cleanall buildall

# Cross compilation
build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} $(GOBUILD) -ldflags "-s -w -X ${PACKAGE_BASE}/version.CommitHash=${GITHASH} -X ${PACKAGE_BASE}/version.BuildDate=${BUILDDATE} -X ${PACKAGE_BASE}/version.Version=${VERSION}" -o ${BUILD_DIR}$(BINARY_UNIX) -v

build-debug:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=amd64 $(GOBUILD) -ldflags "-X ${PACKAGE_BASE}/version.CommitHash=${GITHASH} -X ${PACKAGE_BASE}/version.BuildDate=${BUILDDATE} -X ${PACKAGE_BASE}/version.Version=${VERSION}" -o ${BUILD_DIR}$(BINARY_UNIX) -v

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -rf ${BUILD_DIR}
	rm -rf modules/*

cleanplugins:
	$(GOCLEAN)
	rm -rf modules/*

cleanall: clean cleanplugins

run:
	mkdir ${BUILD_DIR}tmp/
	$(GOBUILD) -a -o ${BUILD_DIR}$(BINARY_NAME) -v ./...
	./${BUILD_DIR}$(BINARY_NAME)

deps:
	$(GO_DEP_FETCH) github.com/prometheus/client_golang/prometheus/promhttp
