#!/bin/bash

# Script to run batch analysis in Docker
# Usage: ./run_analysis.sh <tournament_folder> [workers] [enable_first_round]
# Example: ./run_analysis.sh results_20260210_011606 8 true

TOURNAMENT_FOLDER=${1:-"results_20260210_011606"}
WORKERS=${2:-4}
ENABLE_FIRST_ROUND=${3:-""}

EXTRA_ARGS=""
if [ "$ENABLE_FIRST_ROUND" = "true" ]; then
    EXTRA_ARGS="--enable-first-round-analysis"
fi

echo "Running batch analysis for: $TOURNAMENT_FOLDER"
echo "Workers: $WORKERS"
echo "First round analysis: ${ENABLE_FIRST_ROUND:-false}"
echo "---"

docker run --rm -v "$(pwd)/..:/app" -w /app csgo-abm-analysis \
  python batch_analyze_tournament.py "$TOURNAMENT_FOLDER" --workers "$WORKERS" $EXTRA_ARGS
