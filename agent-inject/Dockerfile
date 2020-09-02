FROM golang:latest as builder
RUN mkdir -p /google-secret-k8s
WORKDIR /google-secret-k8s
COPY . .
RUN make build.binaries

FROM alpine:3.10
WORKDIR /
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' appuser
COPY --from=builder /google-secret-k8s/bin/google-secret-k8s .
COPY --from=builder /google-secret-k8s/cert /cert
USER appuser

CMD ["./gsecret-agent-injector", "server"]
