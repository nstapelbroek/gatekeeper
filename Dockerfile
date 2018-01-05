FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/nstapelbroek/gatekeeper

ENV USER nico
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET -vztTB2-B8TJjmoD

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://nico@localhost:5432/gatekeeper?sslmode=disable

WORKDIR /go/src/github.com/nstapelbroek/gatekeeper

RUN godep go build

EXPOSE 8888
CMD ./gatekeeper