FROM golang:1.11-alpine as build

RUN apk add git make

COPY . /src/feed-collector

WORKDIR /src/feed-collector
RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

COPY --from=build /src/feed-collector/bin/feed-collector /usr/local/bin/feed-collector

HEALTHCHECK CMD curl --silent --show-error --fail http://localhost/health || exit 1
CMD ["/usr/local/bin/feed-collector"]
