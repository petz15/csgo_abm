FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    graphviz \
    libgraphviz-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements file
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy project files (done via volume mount in production)
# COPY . .

# Create output directory
RUN mkdir -p /app/output

# Set environment variables
ENV PYTHONUNBUFFERED=1

# Default command (can be overridden)
CMD ["python", "batch_analyze_tournament.py", "--help"]
