FROM lukechannings/deno

EXPOSE 8080

WORKDIR /app

RUN apt update && apt install -y make

ADD . .

RUN chown -R deno:deno /app

RUN make bundle

ENTRYPOINT []

CMD ["make", "run-bundle"]
