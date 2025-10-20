# Use a small, secure Python base image
FROM python:3.11-slim AS base

# Disable .pyc files and enable unbuffered logs (important for Cloud Run logs)
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Install system dependencies your app might need
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for security
RUN adduser --disabled-password --gecos '' appuser
USER appuser

# Set working directory
WORKDIR /app

# Copy only requirements first (for better caching)
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir -r requirements.txt

# Copy the rest of the app
COPY . .

# Expose the port Cloud Run will use
EXPOSE 8080

# Run the app using Uvicorn, with host 0.0.0.0 and Cloud Run's $PORT variable
# (Cloud Run automatically provides the PORT env var)
CMD exec uvicorn api_endpoint:app --host 0.0.0.0 --port ${PORT:-8000} --proxy-headers
