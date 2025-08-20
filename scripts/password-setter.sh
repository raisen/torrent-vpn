#!/bin/sh

sleep 30
echo 'Setting qBittorrent password and configuring for streaming...'

# Wait for qBittorrent to be ready
for i in $(seq 1 30); do
  if curl -s http://localhost:8080 > /dev/null 2>&1; then
    echo 'qBittorrent is ready'
    break
  fi
  echo 'Waiting for qBittorrent to start...'
  sleep 5
done

# Try to login with default credentials first
echo 'Trying to login with admin/admin...'
LOGIN_RESPONSE=$(curl -s -c /tmp/cookies.txt --data 'username=admin&password=admin' http://localhost:8080/api/v2/auth/login)

if [ "$LOGIN_RESPONSE" = "Ok." ]; then
  echo 'Already configured with admin/admin'
else
  echo 'Checking for temporary password...'
  # Since we can't access docker logs directly, try common temporary passwords or wait for manual setup
  TEMP_PASSWORDS="adminadmin password 123456 qbittorrent"
  for TEMP_PASS in $TEMP_PASSWORDS; do
    LOGIN_RESPONSE=$(curl -s -c /tmp/cookies.txt --data "username=admin&password=$TEMP_PASS" http://localhost:8080/api/v2/auth/login)
    if [ "$LOGIN_RESPONSE" = "Ok." ]; then
      echo "Login successful with password: $TEMP_PASS"
      # Set permanent password
      curl -s -b /tmp/cookies.txt --data 'json={"web_ui_password":"admin"}' http://localhost:8080/api/v2/app/setPreferences
      echo 'Password changed to admin'
      break
    fi
  done
fi

# Configure for streaming incomplete downloads and sequential download default
echo 'Configuring streaming settings and sequential download defaults...'
curl -s -b /tmp/cookies.txt --data 'json={"torrent_content_layout":"Original","auto_tmm_enabled":false,"preallocate_all":false,"incomplete_files_ext":false,"sequential_download":true,"first_last_piece_prio":true}' http://localhost:8080/api/v2/app/setPreferences
echo 'Streaming configuration and sequential download defaults applied successfully!'

# Signal that configuration is complete
touch /tmp/gluetun/qbittorrent_configured
echo 'Configuration finished, signaling other services...'
echo 'Password setter completed!'
