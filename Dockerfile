FROM golang:1.26-alpine AS compiler

WORKDIR /workspace

# Force dependency verification parity
COPY go.mod ./
RUN go mod download

COPY . .

# Compile optimized static binary with debugging overhead completely stripped
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /artifacts/secure-api main.go

FROM alpine:3.19

RUN apk update && apk add --no-cache curl \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /runtime

COPY --from=compiler /artifacts/secure-api .

# Drop absolute execution privileges down to unprivileged workspace user
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./secure-api"]
