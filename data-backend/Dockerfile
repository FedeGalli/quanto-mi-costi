# ---- Base image ----
FROM python:3.11-slim

# Environment setup
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PATH="/usr/local/bin:/home/appuser/.local/bin:$PATH"

# Install build tools (needed for polars)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m appuser
WORKDIR /app

# Copy requirements first (for better Docker layer caching)
COPY requirements.txt .

# Install dependencies as root into /usr/local (so PATH already includes them)
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Copy the rest of the app
COPY . .

# Switch to non-root user
USER appuser

# Expose Cloud Run's port
EXPOSE 8080

# Start the FastAPI app with uvicorn
CMD ["uvicorn", "api_endpoint:app", "--host", "0.0.0.0", "--port", "8000", "--proxy-headers"]
