FROM golang:1.23-bookworm AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /eos-server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /eos-worker ./cmd/worker/main.go

FROM node:22-bookworm-slim AS build-assets-stage

WORKDIR /app
COPY . /app
RUN apt-get update -y && apt-get -y install make
RUN npm i
RUN make build-assets-css
RUN make build-assets-js

FROM debian:bookworm AS server-release-stage

WORKDIR /
COPY --chown=nonroot --from=build-stage /eos-server /eos-server
COPY --chown=nonroot --from=build-stage /app/public /public
EXPOSE 3000
ENTRYPOINT ["/eos-server"]

FROM debian:bookworm AS worker-release-stage

WORKDIR /
COPY --chown=nonroot --from=build-stage /eos-worker /eos-worker
COPY --chown=nonroot --from=build-assets-stage /app/public /public
EXPOSE 3000
ENTRYPOINT ["/eos-worker"]
