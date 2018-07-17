binname=feed-collector
version=devel
prefix=/usr/local

.PHONY: build clean docker fmt test

# Building Commands

build:
	go build -ldflags "-X main.version=$(version)" -o "bin/$(binname)" .

docker:
	env CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$(version)" -o "bin/$(binname).linux" .
	docker build . --tag $(binname)

# Maintenance Commands

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

test:
	go test ./...

clean:
	-rm -rf bin/
