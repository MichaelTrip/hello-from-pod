FROM golang:alpine3.17 AS builder

WORKDIR /src
COPY . .
RUN go build -o dist/

FROM alpine:3.17
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /src/dist/hello-from-pod .
RUN chmod +x /app/hello-from-pod

CMD ["/app/hello-from-pod"]
