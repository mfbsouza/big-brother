# Makefile by Matheus Souza (github.com/mfbsouza)

# project name
PROJECT := big-brother

# paths
BUILDDIR := ./build
DBGDIR   := $(BUILDDIR)/debug
RELDIR   := $(BUILDDIR)/release

ifeq ($(RELEASE),1)
	FLAGS  := '-ldflags "-s -w"'
	BINDIR := $(RELDIR)
else
	FLAGS  :=
	BINDIR := $(DBGDIR)
endif

# targets
.PHONY: all build run format tests clean

all: build

build: $(BINDIR)
	go build ./cmd/web $(FLAGS)

$(BINDIR):
	@mkdir -p $(BINDIR)

run:
	go run ./cmd/web

format:
	gofmt -w -s .

tests:
	echo "todo"

clean:
	rm -rf ./build/
