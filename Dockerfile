FROM golang:1.20-alpine3.17
RUN apk add --no-cache make

WORKDIR /app
ADD . .
RUN make build

CMD bin/service
