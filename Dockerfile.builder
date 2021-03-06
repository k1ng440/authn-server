FROM golang:1.10-alpine
ARG glide_version="v0.13.1"

RUN apk update
RUN apk add --no-cache ca-certificates openssl git make bash gcc musl-dev

RUN wget https://github.com/Masterminds/glide/releases/download/${glide_version}/glide-${glide_version}-linux-amd64.tar.gz
RUN tar -xvzf glide-${glide_version}-linux-amd64.tar.gz
RUN mv linux-amd64/glide /usr/local/bin/glide

RUN go get github.com/benbjohnson/ego/cmd/ego

ENV PATH="${PATH}:/usr/local/go/bin"

WORKDIR /
