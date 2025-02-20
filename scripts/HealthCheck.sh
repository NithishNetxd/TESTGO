#!/bin/bash

# Health check settings
URL="http://localhost:8080/health"  # Use the correct health check endpoint
RETRIES=6  # Total retries (30 seconds total, checking every 5 seconds)
REQUIRED_SUCCESSES=5  # Must get 200 at least 5 times in a row
SUCCESS_COUNT=0

echo "Starting health check..."

for ((i=1; i<=RETRIES; i++)); do
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $URL || echo "000")
  echo "Health check attempt $i: HTTP $HTTP_STATUS"

  if [[ "$HTTP_STATUS" == "200" ]]; then
    ((SUCCESS_COUNT++))
    echo "✅ Success count: $SUCCESS_COUNT/$REQUIRED_SUCCESSES"
    
    if [[ $SUCCESS_COUNT -ge $REQUIRED_SUCCESSES ]]; then
      echo "🎉 Health check passed! Proceeding..."
      exit 0
    fi
  else
    echo "❌ Health check failed. Resetting success count."
    SUCCESS_COUNT=0  # Reset success count if failure occurs
  fi

  sleep 5  # Wait 5 seconds before retrying
done

echo "🚨 Health check failed after multiple attempts! Stopping deployment."
exit 1
