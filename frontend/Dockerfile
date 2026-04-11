FROM oven/bun:alpine AS builder

WORKDIR /app

COPY package.json ./
COPY bun.lock ./
RUN bun install

COPY . .
RUN bun run build-strict

# ==========================================

FROM oven/bun:alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/server.ts ./server.ts

EXPOSE 8080

CMD ["bun", "run", "server.ts"]
