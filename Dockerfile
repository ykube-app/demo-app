# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend (copies the built frontend from stage 1)
FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY backend/ ./
COPY --from=frontend /app/frontend/../backend/static ./static
RUN go build -o server .

# Stage 3: Minimal runtime image
FROM alpine:3.21
COPY --from=backend /app/server /server
ENTRYPOINT ["/server"]