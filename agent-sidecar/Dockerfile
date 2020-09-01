FROM golang:latest as builder
RUN mkdir -p /gsecret-agent
WORKDIR /gsecret-agent
COPY . .
RUN make build.binaries

FROM alpine:3.10
WORKDIR /
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' appuser
COPY --from=builder /gsecret-agent/bin/gsecret-agent .
USER appuser

CMD ["./gsecret-agent"]
