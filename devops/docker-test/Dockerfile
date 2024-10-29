FROM golang:1.23-alpine as go-builder

WORKDIR /app

COPY . .

RUN apk add -u -t build-tools curl git && \
  go build -o server *.go && \
  rm -rf /var/cache

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=go-builder /app/server /app/server

EXPOSE 8080

ENTRYPOINT [ "/app/server" ]
