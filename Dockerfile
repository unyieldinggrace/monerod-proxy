
FROM golang:1.17.7 AS builder

WORKDIR /src
ARG BUILDARGS
COPY . .
RUN go build $BUILDARGS


FROM debian:11.2-slim

WORKDIR /app
COPY --from=builder /src/monerod-proxy ./
ENTRYPOINT ["/app/monerod-proxy"]
