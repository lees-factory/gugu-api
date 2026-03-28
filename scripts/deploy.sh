#!/bin/bash
set -e

APP_DIR="/home/ubuntu/gugu-api"
APP_NAME="gugu-api"

cd "$APP_DIR"

# Stop existing process (PID file 기반)
if [ -f "$APP_DIR/app.pid" ]; then
  OLD_PID=$(cat "$APP_DIR/app.pid")
  if kill -0 "$OLD_PID" 2>/dev/null; then
    kill "$OLD_PID" || true
    sleep 2
  fi
  rm -f "$APP_DIR/app.pid"
fi

# Replace binary
mv "${APP_NAME}.new" "$APP_NAME"
chmod +x "$APP_NAME"

# Start application
nohup ./"$APP_NAME" > "$APP_DIR/app.log" 2>&1 &
echo $! > "$APP_DIR/app.pid"

sleep 2

# Health check
if kill -0 "$(cat "$APP_DIR/app.pid")" 2>/dev/null; then
  echo "Deploy successful - $APP_NAME is running (PID: $(cat "$APP_DIR/app.pid"))"
else
  echo "Deploy failed - $APP_NAME is not running"
  exit 1
fi
