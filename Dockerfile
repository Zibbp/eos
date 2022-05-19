FROM golang:1.18 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o eos cmd/server/main.go

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./eos"]