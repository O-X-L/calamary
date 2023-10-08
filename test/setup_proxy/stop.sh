#!/bin/bash

ruleset="$(nft --handle list ruleset)"

VERSION="$1"

table='default'
table_type='inet'
dnat_chain='prerouting'

# todo: add rule to redirect some ports to tproxy-enabled listener
tproxy='false'
proxy_port="$(cat "/tmp/calamary_${VERSION}.yml" | grep -A2 "mode: 'transparent'" | grep -A1 "tproxy: ${tproxy}" | grep 'port:' | cut -d ' ' -f8)"

redirect_rule="meta l4proto { tcp, udp } ip saddr ${TESTER_HOST} redirect"

if ! grep -q "$redirect_rule" <<< "$ruleset"
then
  exit 0
fi

dnat_chain="$(grep 'type nat hook prerouting' -B1 <<< "$ruleset" | head -n1 | cut -d ' ' -f 2)"
rule_id="$(grep "$redirect_rule" <<< "$ruleset" | cut -d '#' -f 2 | cut -d ' ' -f3)"

echo 'REMOVING NAT RULE'
nft "delete rule ${table_type} ${table} ${dnat_chain} handle ${rule_id}"
