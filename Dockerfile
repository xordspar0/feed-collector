FROM alpine:latest

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

ADD bin/feed-collector.linux /usr/local/bin/feed-collector

CMD /usr/local/bin/feed-collector
