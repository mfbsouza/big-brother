# Makefile by Matheus Souza (github.com/mfbsouza)

# paths
BUILDDIR := ./build
DBGDIR   := $(BUILDDIR)/debug
RELDIR   := $(BUILDDIR)/release

# targets
.PHONY: all build run format tests clean

all: build

build:
	ifeq ($(RELEASE),1)
		go build -ldflags "-s -w" -o $(RELDIR)/
	else
		go build -o $(DBGDIR)/
	endif

run:
	go run .

format:
	gofmt -w -s .

tests:
	echo "todo"

clean:
	rm -rf ./build/
