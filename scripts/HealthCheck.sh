#!/bin/bash

# Health check settings
URL="http://localhost:8080/health"  # Use the correct health check endpoint
RETRIES=6  # Total retries (30 seconds total, checking every 5 seconds)
SUCCESS=0

echo "Starting health check..."

for ((i=1; i<=RETRIES; i++)); do
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $URL || echo "000")
  echo "Health check attempt $i: HTTP $HTTP_STATUS"

  if [[ "$HTTP_STATUS" == "200" ]]; then
    SUCCESS=1
    echo "Health check passed."
    exit 0
  fi

  sleep 5  # Wait 5 seconds before retrying
done

echo "🚨 Health check failed! Stopping deployment."
exit 1
