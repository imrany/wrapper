#!/bin/bash

set -e

echo "Deploying Wrapper..."
docker stop wrapper 2>/dev/null || true
docker rm wrapper 2>/dev/null || true
docker rmi ghcr.io/imrany/wrapper 2>/dev/null || true
docker pull ghcr.io/imrany/wrapper:latest
docker run -d --name wrapper --env-file .env -p 8000:8000 -p 8090:8090 -v ~/.wrapper:/var/opt/wrapper ghcr.io/imrany/wrapper:latest
echo "Deployment complete. Showing logs:"
docker logs wrapper
