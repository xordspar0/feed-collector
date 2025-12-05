binname=feed-collector
version=devel

.PHONY: build clean docker fmt test

# Building Commands

build:
	env CGO_ENABLED=0 go build -ldflags "-X main.version=$(version)" -o "bin/$(binname)" .

docker:
	docker build . --tag $(binname)

# Maintenance Commands

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

test:
	go test ./...

clean:
	-rm -rf bin/
