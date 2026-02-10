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

# Run analysis (results from /mnt/hdd/csgo_abm)
./run_analysis.sh results_20260210_011606 8 false /mnt/hdd/csgo_abm

# Or with first round analysis enabled
./run_analysis.sh results_20260210_011606 8 true /mnt/hdd/csgo_abm
```

**Manual Docker command:**

```bash
docker run --rm \
  -v "$(pwd)/..:/app" \
  -v "/mnt/hdd/csgo_abm:/mnt/results" \
  -w /app \
  csgo-abm-analysis \
  python batch_analyze_tournament.py /mnt/results/results_20260210_011606 --workers 8
```

**Using docker-compose:**

```bash
# Results are auto-mounted from /mnt/hdd/csgo_abm to /mnt/results inside container
docker-compose run --rm csgo-analysis \
  python batch_analyze_tournament.py /mnt/results/results_20260210_011606 --workers 8
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
./run_analysis.sh <tournament_folder> [workers] [enable_first_round] [results_base_path]
```

**Arguments:**
- `tournament_folder`: Name of tournament results folder (required)
- `workers`: Number of parallel workers (default: 4)
- `enable_first_round`: Set to "true" to enable first round analysis (default: false)
- `results_base_path`: Base path where results are stored (default: /mnt/hdd/csgo_abm)

**Examples:**
```bash
# Basic usage with default settings (results from /mnt/hdd/csgo_abm)
./run_analysis.sh results_20260210_011606

# With 8 workers
./run_analysis.sh results_20260210_011606 8

# With first round analysis enabled
./run_analysis.sh results_20260210_011606 8 true

# With custom results location
./run_analysis.sh results_20260210_011606 8 false /mnt/hdd/csgo_abm

# Full example with all parameters
./run_analysis.sh results_20260210_011606 8 true /mnt/hdd/csgo_abm
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
# FVolume Mounts

The Docker setup uses the following volume mounts:

- `..:/app` - Maps parent directory (csgo_abm) to /app for access to scripts and notebooks
- `/mnt/hdd/csgo_abm:/mnt/results` - Maps your results directory to /mnt/results inside container

**Important:** The results location is configured for `/mnt/hdd/csgo_abm` by default. If your results are in a different location, update `docker-compose.yml` or provide the path when using `run_analysis.sh`.

## Notes

- All results and HTML files are saved to your host machine via volume mounting
- The Docker container has access to all files in the parent `csgo_abm` directory
- Jupyter notebooks execute within the container but files are persisted on your host
- The image includes all necessary Python packages for data analysis and visualization
- Results are expected at `/mnt/hdd/csgo_abm` on the host machine (mounted to `/mnt/results` in container)

## Notes

- All results and HTML files are saved to your host machine via volume mounting
- The Docker container has access to all files in the parent `csgo_abm` directory
- Jupyter notebooks execute within the container but files are persisted on your host
- The image includes all necessary Python packages for data analysis and visualization
