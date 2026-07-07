FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o admin-backend .

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget \
  && addgroup -S conelli \
  && adduser -S conelli -G conelli

COPY --from=builder /app/admin-backend /app/admin-backend

USER conelli

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD wget -qO- http://127.0.0.1:8000/health || exit 1

CMD ["/app/admin-backend"]
