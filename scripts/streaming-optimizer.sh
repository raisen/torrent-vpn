#!/bin/sh

echo 'Starting qBittorrent streaming optimizer...'

# Wait for qBittorrent configuration to be complete
while [ ! -f /tmp/gluetun/qbittorrent_configured ]; do
  echo 'Waiting for qBittorrent configuration to complete...'
  sleep 5
done

echo 'qBittorrent is configured! Starting streaming optimization...'
# First run: optimize ALL existing torrents
echo 'Initial optimization of existing torrents...'

if curl -s --data 'username=admin&password=admin' http://localhost:8080/api/v2/auth/login --cookie-jar /tmp/qb.cookie | grep -q 'Ok'; then
  echo 'Logged in successfully, optimizing all video torrents...'
  # Get ALL torrents (not just downloading)
  ALL_TORRENTS=$(curl -s --cookie /tmp/qb.cookie http://localhost:8080/api/v2/torrents/info)
  echo "$ALL_TORRENTS" | jq -r '.[] | select(.name | test("(mkv|mp4|avi|mov|wmv|flv|webm|m4v|S[0-9]|Season|Episode)"; "i")) | select(.seq_dl == false) | .hash' | while read hash; do
    if [ -n "$hash" ]; then
      echo "Enabling sequential download for existing torrent: $hash"
      curl -s --cookie /tmp/qb.cookie --data "hashes=$hash" http://localhost:8080/api/v2/torrents/toggleSequentialDownload
      curl -s --cookie /tmp/qb.cookie --data "hashes=$hash" http://localhost:8080/api/v2/torrents/toggleFirstLastPiecePrio
    fi
  done
fi

# Continuous monitoring for new torrents
while true; do
  sleep 60
  if curl -s --data 'username=admin&password=admin' http://localhost:8080/api/v2/auth/login --cookie-jar /tmp/qb.cookie | grep -q 'Ok'; then
    # Check for new video torrents without sequential download
    TORRENTS=$(curl -s --cookie /tmp/qb.cookie http://localhost:8080/api/v2/torrents/info)
    echo "$TORRENTS" | jq -r '.[] | select(.name | test("(mkv|mp4|avi|mov|wmv|flv|webm|m4v|S[0-9]|Season|Episode)"; "i")) | select(.seq_dl == false) | .hash' | while read hash; do
      if [ -n "$hash" ]; then
        echo "Enabling sequential download for new torrent: $hash"
        curl -s --cookie /tmp/qb.cookie --data "hashes=$hash" http://localhost:8080/api/v2/torrents/toggleSequentialDownload
        curl -s --cookie /tmp/qb.cookie --data "hashes=$hash" http://localhost:8080/api/v2/torrents/toggleFirstLastPiecePrio
      fi
    done
  fi
done
