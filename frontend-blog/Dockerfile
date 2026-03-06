# File: Dockerfile
# Purpose: Build and run the Nuxt public blog as a production SSR container.
# Module: frontend-blog deployment, container image layer.
# Related: nuxt.config.ts, package.json, deploy/docker-compose.yml.

FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

ARG NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
ENV NUXT_PUBLIC_API_BASE=${NUXT_PUBLIC_API_BASE}

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

FROM node:22-alpine

WORKDIR /app

RUN apk add --no-cache libc6-compat && \
    addgroup -S app && \
    adduser -S -G app app

ENV NODE_ENV=production
ENV NITRO_PORT=3000
ENV NITRO_HOST=0.0.0.0

COPY --from=builder /app/.output /app/.output

USER app

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
