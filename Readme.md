A wrapper for gemini ai models

```bash
go run main.go --port=5000 --gemini-api-key=your-key
```

## 🧪 How to Build and Run

### 1. **Build the Docker image**

```bash
docker build -t wrapper .
```

### 2. **Run with flags**

```bash
docker run -p 8000:8000 wrapper --port=8000 --gemini-api-key=your-key-here
```

Or using an environment variable:

```bash
docker run -p 8000:8000 -e GEMINI_API_KEY=your-key-here wrapper --port=8000
```

The app will pick up `GEMINI_API_KEY` via Viper’s `AutomaticEnv()`.

---

## 🧠 Bonus: Use `.env` File in Docker

If you want to use a `.env` file:

### 1. Create `.env`

```env
PORT=8000
GEMINI_API_KEY=your-key-here
```

### 2. Run with env file

```bash
docker run --env-file .env -p 8000:8000 wrapper
```