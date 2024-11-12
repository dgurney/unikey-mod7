FROM golang:1.23 as build
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID -o build
RUN useradd -m -u $UID -g $GID -o -s /bin/bash build
WORKDIR /go/src/unikey-mod7
RUN apt update
RUN apt install -y zip

USER build

# For efficient caching, download deps first
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /go/src/unikey-mod7
ENV GOPATH /go
