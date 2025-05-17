FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o branch-wrangler ./cmd/branch-wrangler

FROM alpine:latest
COPY --from=builder /app/branch-wrangler /usr/local/bin/
ENTRYPOINT ["branch-wrangler"]