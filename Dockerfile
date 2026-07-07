FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o admin-backend .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/admin-backend /app/admin-backend

EXPOSE 8000

CMD ["/app/admin-backend"]
