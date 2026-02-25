FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /shortcuts .

FROM python:3.12-slim
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg ca-certificates && \
    pip install --no-cache-dir yt-dlp && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /shortcuts /app/shortcuts
RUN useradd -r -u 1001 appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/shortcuts"]
