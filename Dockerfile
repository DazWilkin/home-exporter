ARG GOLANG_VERSION=1.24

ARG COMMIT
ARG VERSION

FROM golang:${GOLANG_VERSION} as build

ARG COMMIT
ARG VERSION

WORKDIR /home-exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY collector collector

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/home-exporter \
    ./main.go

FROM gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source=https://github.com/DazWilkin/home-exporter

EXPOSE 9999

ENTRYPOINT ["/home-exporter"]
