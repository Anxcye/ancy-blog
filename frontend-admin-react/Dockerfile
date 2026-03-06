# File: Dockerfile
# Purpose: Build and serve the React admin app as a static production site.
# Module: frontend-admin-react deployment, container image layer.
# Related: vite.config.ts, package.json, deploy/admin/nginx.conf.

FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

FROM nginx:1.29-alpine

COPY deploy/admin/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
