FROM denoland/deno:alpine-1.30.0

EXPOSE 9000

WORKDIR /app

RUN apk add --no-cache make

ADD . .

RUN chown -R deno:deno /app

USER deno

RUN make bundle

CMD make run-bundle
