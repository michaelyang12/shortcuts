# shortcuts

HTTP API that iOS Shortcuts can call to extract information from videos, images, and text using Claude.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check (no auth required) |
| `POST` | `/text` | Analyze text with a prompt |
| `POST` | `/image` | Analyze a base64-encoded image with a prompt |
| `POST` | `/video` | Download a video, extract frames, and analyze with a prompt |

## Setup

### Prerequisites

- Go 1.23+
- [Anthropic API key](https://console.anthropic.com/)
- For `/video`: [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://ffmpeg.org/) installed locally (or use Docker)

### Environment

```bash
cp .env.example .env
```

Fill in `.env`:

```
SHORTCUTS_API_KEY=<any-secret-you-make-up>
ANTHROPIC_API_KEY=<your-anthropic-api-key>
SHORTCUTS_PORT=8080
```

`SHORTCUTS_API_KEY` is a shared secret between your server and your iOS Shortcuts. Pick any random string (e.g. `uuidgen` output).

### Run locally

```bash
source .env && export SHORTCUTS_API_KEY ANTHROPIC_API_KEY SHORTCUTS_PORT
go run .
```

### Run with Docker

```bash
docker compose up --build
```

## Usage

All endpoints (except `/health`) require a `Bearer` token:

```
Authorization: Bearer <your-SHORTCUTS_API_KEY>
```

### POST /text

```bash
curl -s -X POST http://localhost:8080/text \
  -H "Authorization: Bearer $SHORTCUTS_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"text": "The quick brown fox jumps over the lazy dog", "prompt": "Count the words"}'
```

```json
{"result": "The sentence contains 9 words."}
```

### POST /image

```bash
curl -s -X POST http://localhost:8080/image \
  -H "Authorization: Bearer $SHORTCUTS_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "<base64-encoded-image>",
    "media_type": "image/jpeg",
    "prompt": "Describe this image"
  }'
```

`media_type` is optional (defaults to `image/jpeg`). Supported: `image/jpeg`, `image/png`, `image/gif`, `image/webp`.

### POST /video

```bash
curl -s -X POST http://localhost:8080/video \
  -H "Authorization: Bearer $SHORTCUTS_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=...", "prompt": "Summarize this video"}'
```

Downloads the video with yt-dlp, extracts 6 evenly-spaced frames with ffmpeg, and sends them to Claude for analysis. Supports any URL that yt-dlp supports (YouTube, TikTok, etc.).
