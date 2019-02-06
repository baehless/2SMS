#!/bin/bash

set -e # fail if unexpected error

echo "Started installation of 2SMS Scraper application"

# Defaults
INSTALLATION_PATH="$HOME/2SMS/deployment/scraper"
DEFAULT_SOCKET=/run/shm/sciond/default.sock
SERVICE_FILE_NAME=2SMSscraper.service
SERVICE_FILE_LOCATION=/etc/systemd/system
MANAGER_IP='192.33.93.196'
monitoring_download_page="monitoring.scionlab.org/downloads/public/scraper"
IP=$(curl -s ipinfo.io/ip)
[ -f $SC/gen/ia ] && IA=$(cat $SC/gen/ia | sed 's/_/:/g') || { echo "Missing $SC/gen/ia file"; exit 1; }

# Prepare filesystem
echo "Creating installation directory at $INSTALLATION_PATH"
mkdir -p $INSTALLATION_PATH
cd $INSTALLATION_PATH

# TODO: improve download process (check return status, missing files, create a separate function)

echo "Downloading scraper binary"
rm -f scraper
wget "https://$monitoring_download_page/scraper.tar.gz" >/dev/null 2>&1
tar -xzvf scraper.tar.gz
rm -f scraper.tar.gz
chmod +x scraper

echo "Downloading Prometheus"
# TODO

# TODO: make prometheus/prometheus executable

echo "Downloading configuration files"
# TODO


# Check system
echo "Checking SCION installation"
if [ -z $SC ]; then
    echo 'SCION environment variable $SC not set, is SCION installed correctly?'
    exit 1
elif [ ! -d $SC/gen/ISD*/AS*/endhost/certs ]; then
    echo "SCION scraper's cert directory not found. Aborting."
    exit 1
elif [ ! -S $DEFAULT_SOCKET ]; then
    echo "SCION default socket not found at $DEFAULT_SOCKET"
    exit 1
else
    echo "SCION correctly installed"
fi

# Start service
echo "Stopping $SERVICE_FILE_NAME"
sudo systemctl stop $SERVICE_FILE_NAME || true
echo "Reloading Daemons"
sudo systemctl daemon-reload
echo "Starting $SERVICE_FILE_NAME"
sudo systemctl start $SERVICE_FILE_NAME
echo "Enabling $SERVICE_FILE_NAME at boot time"
sudo systemctl enable $SERVICE_FILE_NAME || true # at boot time

# Check service status
if systemctl is-active --quiet 2SMSscraper.service; then
    echo "Successfully installed and started 2SMS Scraper application"
else
    echo 'Failed starting 2SMS Scraper application, something went wrong during installation. Check journal and syslog to have more details.'
    exit 1
fi
