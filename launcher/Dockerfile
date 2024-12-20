FROM node:22-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

FROM base AS node-builder
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build

FROM golang:1.22-bookworm AS go-builder
WORKDIR /app
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o server ./cmd

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=node-builder /app/dist ./dist
COPY --from=go-builder /app/server .

ENV PRIV_KEY_FOLDER=/keys

EXPOSE 3000
CMD ["./server"]
