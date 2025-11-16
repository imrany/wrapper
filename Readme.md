# Wrapper

A lightweight gRPC + RESTful wrapper for AI models. Built in Go,
this service exposes a unified interface for AI generation via both
protocol-level gRPC and HTTP/JSON endpoints. Includes Swagger UI,
Docker support, and flexible configuration via flags or environment variables.

## 🚀 Features

- ✅ gRPC service: `AiService.GenAi(prompt)`
- ✅ RESTful HTTP endpoint: `POST /v1/genai`
- ✅ Swagger UI for interactive API testing
- ✅ Docker-ready with multi-port support
- ✅ Configurable via flags, `.env`, or inline environment variables
- ✅ Graceful shutdown and signal handling
- ✅ Gemini API key integration

---

## 🧪 Quick Start

### 1. **Run Locally**

```bash
go run main.go --port=8000 --api-key=your_key_here --model=gemini-2.0-flash
```

Or use environment variables:

```bash
export PORT=8000
export API_KEY=your_key_here
export MODEL=gemini-2.0-flash
go run main.go
```

### 2. **Build Docker Image**

```bash
docker build -t wrapper .
```

Or pull from GitHub Container Registry:

```bash
docker pull ghcr.io/imrany/wrapper
```

### 3. **Run with Docker**

#### Option A: Inline flags

```bash
docker run -d \
  -p 8000:8000 \  # gRPC
  -p 8090:8090 \  # REST + Swagger
  ghcr.io/imrany/wrapper \
  --port=8000 \
  --api-key=your_key_here
  --model=gemini-2.0-flash
```

#### Option B: `.env` file

Create `.env`:

```env
PORT=8000
API_KEY="your_key_here"
MODEL="gemini-2.0-flash" # or "gpt-5.1-2025-11-13"
```

Run:

```bash
docker run --env-file .env -p 8000:8000 -p 8090:8090 ghcr.io/imrany/wrapper
```

## 📡 API Endpoints

### 1. **REST (HTTP/JSON)**

```http
POST /v1/genai
Content-Type: application/json
Body: "Hello AI"
```

#### Example

```bash
curl -X POST http://localhost:8090/v1/genai \
  -H "Content-Type: application/json" \
  -d '"Hello AI"'
```

### 2. **gRPC**

```proto
service AiService {
  rpc GenAi(GenAiRequest) returns (GenAiResponse);
}
```

#### Example

```bash
grpcurl -insecure localhost:8000 \
  wekalist.api.v1.AiService.GenAi \
  -d '{"prompt": "Hello AI"}'
```

---

### 3. **Swagger UI**

Visit: `http://localhost:8090/swagger/`

## 🧠 Response Format

```json
{
  "prompt": "Hello AI",
  "response": "Hello! How can I help you today?\n",
  "status": null
}
```

## 🔐 Environment Variables

| Variable  | Description      |
| --------- | ---------------- |
| `PORT`    | gRPC server port |
| `API_KEY` | Gemini API key   |

---

## 🛡️ Health Check

```http
GET /healthz
```

Returns `200 OK` with body `ok`.

## 🧠 License & Credits

Built by [Imran](https://github.com/imrany)
Licensed under MIT
