FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o rift-chat-api .

# ==========================================
FROM scratch

WORKDIR /

# Copy the static binary
COPY --from=builder /app/rift-chat-api /rift-chat-api

EXPOSE 8081
EXPOSE 4242/udp

# Direct exec form (no shell needed)
ENTRYPOINT ["/rift-chat-api"]
