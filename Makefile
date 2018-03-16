binname=feed-collector
version=devel
prefix=/usr/local

.PHONY: build clean docker fmt install test uninstall

# Building Commands

build:
	go build -ldflags "-X main.version=$(version)" -o "bin/$(binname)" .

bin/$(binname).linux:
	env GOOS=linux go build -ldflags "-X main.version=$(version)" -o "bin/$(binname).linux" .

docker: bin/$(binname).linux
	docker build . --tag $(binname)

# Installing Commands

install: build squirrelbot.1
	install -Dm 755 "bin/$(binname)" "$(prefix)/bin/$(binname)"
	install -Dm 644 system/squirrelbot.service "$(systemd_unit_path)/squirrelbot.service"

uninstall:
	-rm -f "$(prefix)/bin/$(binname)"
	-rm -f "$(systemd_unit_path)/squirrelbot.service"

# Maintenance Commands

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

test:
	go test ./...

clean:
	-rm -rf bin/
