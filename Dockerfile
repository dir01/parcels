FROM golang:1.20-alpine3.17
RUN apk add --no-cache make
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
ADD . .
RUN CGO_ENABLED=1 make install-dev
RUN make build

CMD bin/service
