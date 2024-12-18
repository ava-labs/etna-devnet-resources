FROM node:22-slim AS node-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM golang:1.22-bookworm AS go-builder
WORKDIR /app
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go build -o server ./cmd

FROM debian:bookworm
WORKDIR /app
COPY --from=node-builder /app/dist ./dist
COPY --from=go-builder /app/server .
EXPOSE 3000
CMD ["./server"]
