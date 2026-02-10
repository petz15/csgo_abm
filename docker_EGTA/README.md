# Docker Setup for EGTA Analysis

This directory contains Docker configuration files to run the EGTA Jupyter notebook and batch tournament analysis in an isolated environment.

## Prerequisites

- Docker installed on your Linux machine
- Docker Compose (optional, but recommended)

## Quick Start

### 1. Build the Docker Image

From the `docker_EGTA` directory:

```bash
cd docker_EGTA
docker build -t csgo-abm-analysis -f Dockerfile ..
```

Or using docker-compose:

```bash
cd docker_EGTA
docker-compose build
```

### 2. Run Batch Analysis

**Using the helper script (recommended):**

```bash
# Make the script executable (first time only)
chmod +x run_analysis.sh

# Run analysis
./run_analysis.sh results_20260210_011606 8 true
```

**Manual Docker command:**

```bash
docker run --rm -v "$(pwd)/..:/app" -w /app csgo-abm-analysis \
  python batch_analyze_tournament.py results_20260210_011606 --workers 8
```

**Using docker-compose:**

```bash
docker-compose run --rm csgo-analysis \
  python batch_analyze_tournament.py results_20260210_011606 --workers 8
```

### 3. Run Jupyter Notebook

**Using the helper script (recommended):**

```bash
# Make the script executable (first time only)
chmod +x run_jupyter.sh

# Start Jupyter server
./run_jupyter.sh 8888
```

**Manual Docker command:**

```bash
docker run --rm -p 8888:8888 -v "$(pwd)/..:/app" -w /app csgo-abm-analysis \
  jupyter notebook --ip=0.0.0.0 --port=8888 --no-browser --allow-root
```

**Using docker-compose:**

```bash
docker-compose run --rm -p 8888:8888 csgo-analysis \
  jupyter notebook --ip=0.0.0.0 --port=8888 --no-browser --allow-root
```

## Helper Scripts

### run_analysis.sh

Runs batch tournament analysis with customizable parameters.

**Syntax:**
```bash
./run_analysis.sh <tournament_folder> [workers] [enable_first_round]
```

**Arguments:**
- `tournament_folder`: Path to tournament results folder (required)
- `workers`: Number of parallel workers (default: 4)
- `enable_first_round`: Set to "true" to enable first round analysis (default: false)

**Examples:**
```bash
# Basic usage with default workers
./run_analysis.sh results_20260210_011606

# With 8 workers
./run_analysis.sh results_20260210_011606 8

# With first round analysis enabled
./run_analysis.sh results_20260210_011606 8 true
```

### run_jupyter.sh

Starts Jupyter notebook server for manual notebook execution.

**Syntax:**
```bash
./run_jupyter.sh [port]
```

**Arguments:**
- `port`: Port number for Jupyter server (default: 8888)

**Examples:**
```bash
# Start on default port 8888
./run_jupyter.sh

# Start on custom port
./run_jupyter.sh 9999
```

## File Structure

```
docker_EGTA/
├── Dockerfile              # Docker image definition
├── requirements.txt        # Python dependencies
├── docker-compose.yml      # Docker Compose configuration
├── .dockerignore          # Files to exclude from Docker context
├── run_analysis.sh        # Helper script for batch analysis
├── run_jupyter.sh         # Helper script for Jupyter server
└── README.md              # This file
```

## Common Tasks

### Update Python Dependencies

Edit `requirements.txt` and rebuild:

```bash
docker build -t csgo-abm-analysis -f Dockerfile ..
```

### Clean Up Docker Resources

```bash
# Remove containers
docker-compose down

# Remove image
docker rmi csgo-abm-analysis

# Clean up all unused Docker resources
docker system prune -a
```

### Troubleshooting

**Permission Issues:**
If you encounter permission errors with output files:
```bash
# Run with current user ID
docker run --rm -u $(id -u):$(id -g) -v "$(pwd)/..:/app" -w /app csgo-abm-analysis \
  python batch_analyze_tournament.py results_20260210_011606
```

**Container Won't Start:**
Check Docker logs:
```bash
docker logs csgo-abm-analysis
```

**Port Already in Use:**
Change the port mapping or kill the process using the port:
```bash
# Find process using port 8888
lsof -i :8888

# Use different port
./run_jupyter.sh 9999
```

## Notes

- All results and HTML files are saved to your host machine via volume mounting
- The Docker container has access to all files in the parent `csgo_abm` directory
- Jupyter notebooks execute within the container but files are persisted on your host
- The image includes all necessary Python packages for data analysis and visualization
