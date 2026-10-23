FROM golang:1.27-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \ 
    -ldflags="-s -w" \
    -o linebk-api ./cmd/http
    
FROM gcr.io/distroless/static-debian13:nonroot AS runner

WORKDIR /app

COPY --from=builder --chown=65532:65532 /app/linebk-api /app/linebk-api
COPY --from=builder --chown=65532:65532 /app/config /app/config

ENV PORT=8080
EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app/linebk-api"]