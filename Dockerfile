FROM golang:1.13 AS build-env
# GOPATH is /go
WORKDIR  /go/src/github.com/nstapelbroek/gatekeeper

COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o gatekeeper ./cmd/gatekeeper

FROM alpine:3.10

ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/nstapelbroek/gatekeeper"

RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/nstapelbroek/gatekeeper/gatekeeper /

ENV APP_ENV=release
ENTRYPOINT ["/gatekeeper"]