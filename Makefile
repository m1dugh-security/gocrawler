CC=go build
LDFLAGS=

TARGET=gocrawler
MAIN=./cmd/gocrawler/

.PHONY: all build

all: build

build:
	$(CC) -ldflags '$(LDFLAGS)' -o $(TARGET) '$(MAIN)'

clean:
	rm -f $(TARGET)

tidy:
	go mod tidy
