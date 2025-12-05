# feed-collector

A simple server that aggregates feeds and returns the number of unread articles
in each

## Building

```sh
make
```

## Running

feed-collector is configured using command line flags:

```sh
feed-collector --port 8080 --debug
```

You can also configure feed-collector's options using environment variables:

```sh
FEED_COLLECTOR_PORT=8080 FEED_COLLECTOR_DEBUG=true feed-collector
```

### General Options

```
--port value, -p value           The port to run the server on (default: 80) [$FEED_COLLECTOR_PORT]
--debug                          Enable debug messages [$FEED_COLLECTOR_DEBUG]
--help, -h                       Show this help message
--version, -v                    Show the version
```

### Feed Options

If you run feed-collector without any options, it will not know where to collect
any feeds from. You should specify at least one feed to collect by using the
following options:

```
--nextcloud-news-host value      The hostname of the Nextclout News Server [$NEXTCLOUD_NEWS_HOST]
--nextcloud-news-user value      The user to use for accessing Nextcloud News [$NEXTCLOUD_NEWS_USER]
--nextcloud-news-password value  The password to use for accessing Nextcloud News [$NEXTCLOUD_NEWS_PASSWORD]
```

### Docker

feed-collector comes with a Dockerfile, which can be used to build the app as a
Docker image. To build the Docker image, run

```sh
make docker
```

This will build a Docker image tagged as "feed-collector". To run it as a Docker
container, specify feed-collector's options as environment vaiables:

```sh
docker run -e NEXTCLOUD_NEWS_HOST=https://exmaple.com -e NEXTCLOUD_NEWS_USER=me [-e ...] -d feed-collector
```

You can also use Docker Compose instead of using `docker run` directly. For more
information about Docker Compose, read
[the docs](https://docs.docker.com/compose/). Here is an example
docker-compose.yml file:

```yaml
version: '3'

services:
  feed-collector:
    image: feed-collector
    ports:
      - 80:80
    restart: always
    environment:
      - NEXTCLOUD_NEWS_HOST=https://example.com
      - NEXTCLOUD_NEWS_USER
      - NEXTCLOUD_NEWS_PASSWORD
```

Note that if you don't specify a value for an environment variable in a
compose file, it will use values from the host environment. You could populate
environment variables like this:

```
env NEXTCLOUD_NEEWS_USER=me NEXTCLOUD_NEWS_PASSWORD=asdf docker-compose up -d
```

## Copying

feed-collector is released under the terms of the GNU GPL version 3 or later.
For more information, see the LICENSE file.
