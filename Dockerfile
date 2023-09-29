FROM golang:1.21-alpine
RUN apk add --no-cache make gcc musl-dev

WORKDIR /app
ADD . .
RUN CGO_ENABLED=1 make install-dev
RUN make build

CMD bin/service
