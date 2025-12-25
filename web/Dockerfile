# ===================================================
# Base stage
# ===================================================
FROM oven/bun:1.3.0-alpine AS base
RUN apk add --no-cache ca-certificates curl
WORKDIR /app

# ===================================================
# Dependencies for build stage (includes devDeps)
# ===================================================
FROM base AS deps

ENV HUSKY=0

COPY bun.lock* package.json ./

# Install all dependencies for build (dev + prod)
RUN bun install --frozen-lockfile

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

ARG APP_NAME="Masjid Muslim"
ENV APP_NAME=$APP_NAME

# Build SvelteKit for production
RUN bun -b run build

# ===================================================
# Dependencies for production (only prodDeps)
# ===================================================
FROM base AS prod-deps

ENV HUSKY=0

COPY bun.lock* package.json ./

# Install only production deps
RUN bun install --production --frozen-lockfile

# ===================================================
# Runner stage (minimal, production)
# ===================================================
FROM base AS runner
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
CMD ["bun", "./build/index.js"]
