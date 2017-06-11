# Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VET_REPORT = vet.report
TEST_REPORT = tests.xml

BINARY = dockyard
PKG = github.com/Huawei/containerops/dockyard
REGISTRY = docker.io/containerops/dockyard
GOARCH = amd64

VERSION? = ?
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

BUILD_DIR=${GOPATH}/src/github.com/Huawei/containerops/${BINARY}
CURRENT_DIR=$(shell pwd)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: clean test vet linux darwin windows

linux: 
		cd ${BUILD_DIR}; \
		GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; \
		cd - >/dev/null ; \
		mv ${BINARY}-linux-${GOARCH} ${BUILD_DIR}/bin

darwin:
		cd ${BUILD_DIR}; \
		GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \
	  cd - >/dev/null; \
		mv ${BINARY}-darwin-${GOARCH} ${BUILD_DIR}/bin

windows:
		cd ${BUILD_DIR}; \
	  GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \
	  cd - >/dev/null; \
		mv ${BINARY}-windows-${GOARCH}.exe ${BUILD_DIR}/bin

test:
		if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
		cd ${BUILD_DIR}; \
		godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \
	  cd - >/dev/null; \
		mv ${TEST_REPORT} ${BUILD_DIR}/bin

vet:
		cd ${BUILD_DIR}; \
		godep go vet ./... > ${VET_REPORT} 2>&1 ; \
		cd - >/dev/null; \
		mv ${VET_REPORT} ${BUILD_DIR}/bin

fmt:
		cd ${BUILD_DIR}; \
		go fmt $$(go list ./... | grep -v /vendor/) ; \
	  cd - >/dev/null

clean:
		rm -f ${BUILD_DIR}/bin/${TEST_REPORT}
		rm -f ${BUILD_DIR}/bin/${VET_REPORT}
		rm -f ${BUILD_DIR}/bin/${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean
