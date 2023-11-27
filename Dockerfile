FROM golang:1.21-alpine
RUN apk add --no-cache make gcc musl-dev

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN CGO_ENABLED=1 make install-dev

ADD . .
RUN make build

CMD bin/service
