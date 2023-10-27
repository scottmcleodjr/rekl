APPNAME  = rekl
SRC      = $(wildcard *.go)
BINDIR   = bin

.PHONY: all build run beep clean

all: build

build: $(BINDIR) $(BINDIR)/$(APPNAME)

run: build
	$(BINDIR)/$(APPNAME)

beep: build
	$(BINDIR)/$(APPNAME) -beep

clean:
	rm -rf ${BINDIR}

$(BINDIR):
	mkdir -p $@

$(BINDIR)/$(APPNAME): $(BINDIR) $(SRC) go.mod go.sum
	go build -o $@ $(SRC)