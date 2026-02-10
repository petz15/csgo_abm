#!/bin/bash

# Script to run Jupyter notebook server in Docker
# Usage: ./run_jupyter.sh [port]
# Example: ./run_jupyter.sh 8888

PORT=${1:-8888}

echo "Starting Jupyter notebook server on port $PORT"
echo "Access at: http://localhost:$PORT"
echo "---"

docker run --rm -p "$PORT:8888" -v "$(pwd)/..:/app" -w /app csgo-abm-analysis \
  jupyter notebook --ip=0.0.0.0 --port=8888 --no-browser --allow-root
