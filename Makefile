SHELL = /bin/sh

app_name  = rekl
main_file = app.go
src_files = $(wildcard *.go) $(wildcard */*.go) go.mod go.sum
bin_dir   = bin

all: build
.PHONY: all

build: $(bin_dir) $(bin_dir)/$(app_name)
.PHONY: build

run: build
	$(bin_dir)/$(app_name)
.PHONY: run

beep: build
	$(bin_dir)/$(app_name) -beep
.PHONY: beep

clean:
	rm -rf ${bin_dir}
.PHONY: clean

test:
	go test -v -cover -race ./...
.PHONY: test

$(bin_dir):
	mkdir -p $@

$(bin_dir)/$(app_name): $(bin_dir) $(src_files)
	go build -o $@ $(main_file)
