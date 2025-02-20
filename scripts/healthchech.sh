#!/bin/bash

# Health check settings
URL="http://localhost:8080/health"  # Change to your actual health check endpoint
RETRIES=6  # Total retries (30 seconds total, checking every 5 seconds)
SUCCESS=0

echo "Starting health check..."

for ((i=1; i<=RETRIES; i++)); do
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $URL || echo "000")
  echo "Health check attempt $i: HTTP $HTTP_STATUS"

  if [[ "$HTTP_STATUS" == "200" ]]; then
    SUCCESS=1
    break
  fi

  sleep 5  # Wait 5 seconds before retrying
done

if [[ $SUCCESS -eq 1 ]]; then
  echo "Health check passed. Proceeding..."
  exit 0
else
  echo "Health check failed! Stopping deployment."
  exit 1
fi
