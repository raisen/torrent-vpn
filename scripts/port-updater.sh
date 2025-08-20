#!/bin/sh

echo 'Port updater started, waiting for qBittorrent configuration...'

# Wait for password setter to create the signal file
while [ ! -f /tmp/gluetun/qbittorrent_configured ]; do
  echo 'Waiting for password setter to finish configuration...'
  sleep 5
done

echo 'qBittorrent is configured! Starting port forwarding automation...'
while true; do
  if [ -f /tmp/gluetun/forwarded_port ]; then
    PORT=$(cat /tmp/gluetun/forwarded_port)
    echo "Updating qBittorrent port to: $PORT"
    curl -s --data 'username=admin&password=admin' http://localhost:8080/api/v2/auth/login > /tmp/cookie.txt
    COOKIE=$(grep SID /tmp/cookie.txt | cut -d= -f2)
    if [ -n "$COOKIE" ]; then
      curl -s -H "Cookie: SID=$COOKIE" --data "json={\"listen_port\":$PORT,\"upnp\":false}" http://localhost:8080/api/v2/app/setPreferences
      echo " - Port updated successfully to $PORT"
    else
      echo " - Failed to login to qBittorrent"
    fi
  else
    echo 'Waiting for VPN to establish port forwarding...'
  fi
  sleep 60
done
