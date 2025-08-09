A wrapper for gemini ai models

```bash
go run main.go --port=5000 --api-key=your-key
```

## 🧪 How to Build and Run

### 1. **Build the Docker image**

```bash
docker build -t wrapper .
```
or Pull the image for github
```bash
docker pull ghcr.io/imrany/wrapper
```

### 2. **Run with flags**

```bash
docker run -d \
  -p 8000:8000 \  # gRPC
  -p 8090:8090 \  # REST + Swagger
  ghcr.io/imrany/wrapper \
  --port=8000 \
  --api-key=your_key_here
```

Or using an environment variable:

```bash
docker run -d \
  -p 8000:8000 \
  -p 8090:8090 \
  --env-file .env \
  ghcr.io/imrany/wrapper \
  --port=8000 \
  --api-key=your_key_here
```
or inline env

```bash
docker run -d \
  -p 8000:8000 \
  -p 8090:8090 \
  -e API_KEY=your_key_here \
  ghcr.io/imrany/wrapper
```

The app will pick up `API_KEY` via Viper’s `AutomaticEnv()`.

---

## 🧠 Bonus: Use `.env` File in Docker

If you want to use a `.env` file:

### 1. Create `.env`

```env
PORT=8000
API_KEY=your-key-here
```

### 2. Run with env file

```bash
docker run --env-file .env -p 8000:8000 wrapper
```

```bash
curl -X POST http://localhost:8090/api/v1/genai \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello"}'
```