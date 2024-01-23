# Build Geth in a stock Go builder container
FROM golang:1.20-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-ethereum/
COPY go.sum /go-ethereum/
RUN cd /go-ethereum && go mod download

ADD . /go-ethereum
RUN cd /go-ethereum && go run build/ci.go install -static

# Pull all binaries into a second stage deploy alpine container
FROM ubuntu:latest

RUN echo 'Etc/UTC' > /etc/timezone && \
    echo 'TZif2UTCTZif2UTC\nUTC0' > /etc/localtime && \
    apt-get update && \
    apt-get install -y ca-certificates ntp \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*


COPY --from=builder /go-ethereum/build/bin/* /usr/local/bin/

RUN mkdir -p /app/ \
    && groupadd -g 1000 app \
    && useradd -r -s /bin/sh -d /app/ -u 1000 -g 1000 app \
    && chown -R app:app /app/ \
    && apt update \
    && apt install -y net-tools iputils-ping  \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /app/
USER app

EXPOSE 8545 8546 30303 30303/udp
