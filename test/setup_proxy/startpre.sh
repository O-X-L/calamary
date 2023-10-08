#!/bin/bash

ruleset="$(sudo nft --handle list ruleset)"

VERSION="$1"

table='default'
table_type='inet'
dnat_chain='prerouting'

# todo: add rule to redirect some ports to tproxy-enabled listener
tproxy='false'
proxy_port="$(cat "/tmp/calamary_${VERSION}.yml" | grep -A2 "mode: 'transparent'" | grep -A1 "tproxy: ${tproxy}" | grep 'port:' | cut -d ' ' -f8)"

if grep -q "meta l4proto { tcp, udp } ip saddr ${TESTER_HOST} redirect" <<< "$ruleset"
then
  echo "ERROR: NAT RULE EXISTS!"
  exit 1
fi

if ! grep -q 'type nat hook prerouting' <<< "$ruleset"
then
  sudo nft "add chain ${table_type} ${table} ${dnat_chain} { type nat hook prerouting priority -100; }"
else
  dnat_chain="$(grep 'type nat hook prerouting' -B1 <<< "$ruleset" | head -n1 | cut -d ' ' -f 2)"
fi

echo 'ADDING NAT RULE'
sudo nft "add rule ${table_type} ${table} ${dnat_chain} meta l4proto { tcp, udp } ip saddr ${TESTER_HOST} redirect to ${proxy_port}"
