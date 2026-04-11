FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rift-chat-api .

# ==========================================
FROM alpine:latest  

# Add CA certificates so the Go server can make secure HTTPS requests if needed
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/rift-chat-api .
RUN chmod +x ./rift-chat-api

EXPOSE 8081
EXPOSE 4242/udp

CMD ["./rift-chat-api"]
