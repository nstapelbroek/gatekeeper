FROM alpine:3.7

ADD https://github.com/just-containers/s6-overlay/releases/download/v1.21.2.2/s6-overlay-amd64.tar.gz /tmp/
COPY ./dev/docker/etc /etc
ENTRYPOINT ["/init"]

# Run container setup commands
RUN apk add --no-cache tar ca-certificates \
    && mkdir /app \
    && adduser -S golang\
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C /

COPY gatekeeper /app/gatekeeper

EXPOSE 8080