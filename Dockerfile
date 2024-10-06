FROM golang:1.23-bookworm AS build-server-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /avalon-server ./cmd/server/main.go

FROM node:22-bookworm-slim AS build-assets-stage

WORKDIR /app
COPY . /app
RUN apt-get update -y && apt-get -y install make
RUN make build-assets-css
RUN make build-assets-js

FROM debian:bookworm AS release-stage

WORKDIR /
COPY --chown=nonroot --from=build-server-stage /avalon-server /avalon-server
COPY --chown=nonroot --from=build-assets-stage /app/public /public
EXPOSE 3000
ENTRYPOINT ["/avalon-server"]