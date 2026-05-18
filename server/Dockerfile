# ---- Build stage ----
FROM golang:1.22-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /server /app/server
COPY config.yaml /app/config.yaml
COPY migrations /app/migrations

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["/app/server"]
