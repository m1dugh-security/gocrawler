CC=go build
LDFLAGS=

TARGET=crawler
MAIN=cmd/crawler/main.go

.PHONY: all build

all: build

build:
	$(CC) -ldflags '$(LDFLAGS)' -o $(TARGET) '$(MAIN)'

clean:
	rm -f $(TARGET)

tidy:
	go mod tidy
