FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

ADD bin/feed-collector.linux /usr/local/bin/feed-collector

HEALTHCHECK CMD curl -f http://localhost/health || exit 1
CMD ["/usr/local/bin/feed-collector"]
