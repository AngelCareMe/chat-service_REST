# test.Dockerfile (временно создай в корне)
FROM alpine:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN ls -la