#!/bin/bash
set -e

APP_DIR="/home/ubuntu/gugu-api"
APP_NAME="gugu-api"

cd "$APP_DIR"

# Stop existing process
if pgrep -f "$APP_NAME" > /dev/null; then
  pkill -f "$APP_NAME" || true
  sleep 2
fi

# Replace binary
mv "${APP_NAME}.new" "$APP_NAME"
chmod +x "$APP_NAME"

# Start application
nohup ./"$APP_NAME" > "$APP_DIR/app.log" 2>&1 &

sleep 2

# Health check
if pgrep -f "$APP_NAME" > /dev/null; then
  echo "Deploy successful - $APP_NAME is running"
else
  echo "Deploy failed - $APP_NAME is not running"
  exit 1
fi
