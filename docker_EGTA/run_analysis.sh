#!/bin/bash

# Script to run batch analysis in Docker
# Usage: ./run_analysis.sh <tournament_folder> [workers] [enable_first_round] [results_base_path]
# Example: ./run_analysis.sh results_20260210_011606 8 true /mnt/hdd/csgo_abm

TOURNAMENT_FOLDER=${1:-"results_20260210_011606"}
WORKERS=${2:-4}
ENABLE_FIRST_ROUND=${3:-""}
RESULTS_BASE_PATH=${4:-"/mnt/hdd/csgo_abm"}

EXTRA_ARGS=""
if [ "$ENABLE_FIRST_ROUND" = "true" ]; then
    EXTRA_ARGS="--enable-first-round-analysis"
fi

echo "Running batch analysis for: $TOURNAMENT_FOLDER"
echo "Workers: $WORKERS"
echo "First round analysis: ${ENABLE_FIRST_ROUND:-false}"
echo "Results location: $RESULTS_BASE_PATH"
echo "---"

sudo docker run --rm \
  -v "$(pwd)/..:/app" \
  -v "$RESULTS_BASE_PATH:/mnt/results" \
  -w /app \
  csgo-abm-analysis \
  python batch_analyze_tournament.py "/mnt/results/$TOURNAMENT_FOLDER" --workers "$WORKERS" $EXTRA_ARGS
