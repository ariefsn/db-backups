# ===================================================
# Base stage
# ===================================================
# Install/build stages are pinned to the BUILD platform (the runner's native
# arch) so npm never runs under QEMU emulation when producing arm64 images.
# Only the final `runner` stage is built for the TARGET platform — it just
# copies the arch-independent SvelteKit build + pure-JS node_modules.
# FROM oven/bun:1.3.5-alpine AS base
FROM --platform=$BUILDPLATFORM node:22-alpine AS base
RUN apk add --no-cache ca-certificates curl
WORKDIR /app

# ===================================================
# Dependencies for build stage (includes devDeps)
# ===================================================
FROM base AS deps

ENV HUSKY=0

# Use the lockfile for fast, reproducible installs (npm ci).
COPY package.json package-lock.json ./
RUN npm ci

# ===================================================
# Builder stage
# ===================================================
FROM base AS builder
WORKDIR /app

# Copy deps with devDeps
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/package.json ./package.json
COPY . .

ENV NODE_ENV=production

# Optional: increase heap memory if build fails for large projects
ENV NODE_OPTIONS=--max-old-space-size=8192

ARG APP_NAME="Database Backup Manager"
ENV APP_NAME=$APP_NAME

# Build SvelteKit for production
RUN npm run build

# ===================================================
# Dependencies for production (only prodDeps)
# ===================================================
FROM base AS prod-deps

ENV HUSKY=0

# Install only production deps from the lockfile.
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# ===================================================
# Runner stage (minimal, production) — built for the TARGET platform
# ===================================================
FROM node:22-alpine AS runner
RUN apk add --no-cache ca-certificates curl
WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy build output and production deps
COPY --from=builder /app/build ./build
COPY --from=prod-deps /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Environment
ENV NODE_ENV=production
ARG PORT=5173
ENV PORT=$PORT
EXPOSE $PORT

# Run app
# CMD ["bun", "./build/index.js"]
CMD ["node", "./build/index.js"]
