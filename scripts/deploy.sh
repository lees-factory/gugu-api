#!/bin/bash
set -e

APP_DIR="/home/ubuntu/gugu-api"
APP_NAME="gugu-api"

cd "$APP_DIR"

# Replace binary
mv "${APP_NAME}.new" "$APP_NAME"
chmod +x "$APP_NAME"

# Restart via systemd
sudo systemctl restart gugu-api

sleep 2

# Health check
if systemctl is-active --quiet gugu-api; then
  echo "Deploy successful - gugu-api is running"
else
  echo "Deploy failed - gugu-api is not running"
  sudo journalctl -u gugu-api --no-pager -n 20
  exit 1
fi
