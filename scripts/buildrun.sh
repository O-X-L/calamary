#!/bin/bash

set -e

PROXY_USER='proxy'
CERT_CN="/C=AT/ST=Styria/L=/O=/OU=/CN=Calamary Forward Proxy"

binary="/tmp/calamary_$(date +"%s")"
echo "TEST: Binary: '${binary}'"


if ! [ -f '/tmp/calamary.crt' ]
then
  echo ''
  echo 'GENERATING DUMMY CERTS...'
  echo ''
  openssl req -x509 -newkey rsa:4096 -keyout /tmp/calamary.key -out /tmp/calamary.crt -sha256 -days 60 -nodes -subj "$CERT_CN"
fi

echo ''
echo "EXECUTING CALAMARY AS USER ${PROXY_USER}"
echo ''
cd "$(dirname "$0")/../lib/main"
go build -o "$binary"
chmod +x "$binary"
sudo su "$PROXY_USER" --shell /bin/bash -c "$binary"
