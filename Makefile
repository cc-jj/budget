.PHONY: build clean dev

build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin

dev:
	go tool pulse -c=pulse.json