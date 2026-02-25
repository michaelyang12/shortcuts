# shortcuts — iOS Shortcuts API Server

## Overview

A Go HTTP API that takes media input + a natural language prompt, processes it through Claude's vision/text API, and returns the result. Designed to be called from iOS Shortcuts via Cloudflare Tunnel on a home server.

```
iOS Shortcut → Cloudflare Tunnel → shortcuts API → [preprocess input] → Claude API → response
```

## Endpoints

All endpoints accept JSON and return `{result: string}` or `{error: string}`.

| Endpoint | Input | Preprocessing |
|----------|-------|---------------|
| `POST /video` | `{url, prompt}` | yt-dlp downloads video → ffmpeg extracts ~6 evenly-spaced frames → sends frames + prompt to Claude vision |
| `POST /image` | `{image (base64), prompt}` | Sends image + prompt to Claude vision |
| `POST /text` | `{text, prompt}` | Sends text + prompt to Claude (text-only) |

## Auth

Static API key via `Authorization: Bearer <key>` header, configured in `SHORTCUTS_API_KEY` env var.

## Project Structure

```
shortcuts/
├── main.go              # entry point, router setup, middleware
├── handler/
│   ├── video.go         # POST /video
│   ├── image.go         # POST /image
│   └── text.go          # POST /text
├── claude/
│   └── client.go        # Claude API client (text + vision)
├── media/
│   └── video.go         # yt-dlp download + ffmpeg frame extraction
├── middleware/
│   └── auth.go          # API key auth
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── .env.example
```

## Tech Choices

- **Language:** Go — good HTTP ergonomics, easy exec for yt-dlp/ffmpeg, simple Docker builds
- **LLM:** Anthropic Claude only (direct API calls, no conduit/Fabric dependency)
- **Prompts:** Custom-written, borrowing from Fabric patterns only when a good fit exists
- **Dependencies:** Go standard library for HTTP. yt-dlp + ffmpeg as system binaries in Docker image.
- **Video frames:** ~6 evenly-spaced frames per video

## Env Vars

- `SHORTCUTS_API_KEY` — static Bearer token for auth
- `ANTHROPIC_API_KEY` — Claude API key
- `SHORTCUTS_PORT` — server port (default 8080)

## Deployment

Docker multi-stage build: compile Go binary in build stage, copy into slim image with yt-dlp and ffmpeg installed. Docker Compose for orchestration. Exposed via Cloudflare Tunnel.
