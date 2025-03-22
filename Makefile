# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

APP=rhilex

distro=$(shell grep -oP '(?<=^PRETTY_NAME=")([^"]+)' /etc/os-release)
kernel=$(shell uname -r)
host=$(shell hostname)
ip=$(shell hostname -I | awk '{print $$1}')
memory=$(shell free -m | awk 'NR==2{printf "%.2fGB\n", $$2/1000}')
disk=$(shell df -h | awk '$$NF=="/"{printf "%s\n", $$2}')
arch=$(shell uname -m)

VERSION := $(shell git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD)
HASH := $(shell git rev-parse --short HEAD)

XVersion=-X 'github.com/hootrhino/rhilex/typex.MainVersion=$(VERSION)-${HASH}'
FLAGS="$(XVersion) -s -w -linkmode external -extldflags -static"

GO_BUILD_OPTIONS = -trimpath -ldflags $(FLAGS)

.PHONY: all
all: info build

info:
	@echo "\e[41m[*] Distro \e[0m: \e[36m${distro}\e[0m"
	@echo "\e[41m[*] Arch   \e[0m: \e[36m${arch}\e[0m"
	@echo "\e[41m[*] Kernel \e[0m: \e[36m${kernel}\e[0m"
	@echo "\e[41m[*] Memory \e[0m: \e[36m${memory}\e[0m"
	@echo "\e[41m[*] Host   \e[0m: \e[36m${host}\e[0m"
	@echo "\e[41m[*] IP     \e[0m: \e[36m${ip}\e[0m"
	@echo "\e[41m[*] Disk   \e[0m: \e[36m${disk}\e[0m"

build:
	CGO_ENABLED=1 GOOS=linux go generate
	CGO_ENABLED=1 GOOS=linux go build $(GO_BUILD_OPTIONS) -o ${APP}

.PHONY: x64linux
x64linux:
	go generate
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-x64linux

.PHONY: windows
windows:
	go generate
	# sudo apt install gcc-mingw-w64-x86-64	gcc-mingw-w64
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-windows.exe

.PHONY: arm32
arm32:
	go generate
	# sudo apt install gcc-arm-linux-gnueabi
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-arm32linux

.PHONY: arm64
arm64:
	go generate
	# sudo apt install gcc-aarch64-linux-gnu
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-arm64linux

.PHONY: mips32
mips32:
	go generate
	# sudo apt install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips CGO_ENABLED=1 CC=mips-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-mips32linux

.PHONY: mips64
mips64:
	go generate
	# sudo apt install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips64 CGO_ENABLED=1 CC=mips-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-mips64linux

.PHONY: mips32le
mips32le:
	go generate
	# sudo apt install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mipsle CGO_ENABLED=1 GOMIPS=softfloat CC=mipsel-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-mipslelinux

.PHONY: mips64le
mips64le:
	go generate
	# sudo apt install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips64le CGO_ENABLED=1 GOMIPS=softfloat CC=mipsel-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-mipslelinux

.PHONY: riscv64
riscv64:
	go generate
	# sudo apt install g++-riscv64-linux-gnu gcc-riscv64-linux-gnu
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=1 CC=riscv64-linux-gnu-gcc go build $(GO_BUILD_OPTIONS) -o ${APP}-riscv64linux

.PHONY: release
release:
	bash ./release_pkg.sh

.PHONY: run
run:
	go run -race run

.PHONY: test
test:
	go test ${APP}/test -v

.PHONY: cover
cover:
	go test ${APP}/test -v -cover

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: clean
clean:
	go clean
	rm -rf _release zold zbackup zupgrade* \
	    ${APP}-arm* ${APP}-riscv* ${APP}-mips* \
		license.* \
		*.exe *.db *.txt *.txt.gz .encrypt.p* rhilex.pid \

