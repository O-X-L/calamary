#!/bin/bash

TEST_TARGET='135.181.170.219'
PROXY_UID=13  # exclude traffic from proxy itself to keep us from having a loop

sudo iptables -t nat -I OUTPUT -d "$TEST_TARGET" -p tcp -m owner ! --uid-owner "$PROXY_UID" -j DNAT --to-destination 127.0.0.1:4128

sudo iptables -L -t nat

# to delete:
#  sudo iptables -t nat -L --line-numbers
#  sudo iptables -t nat -D OUTPUT ${NR}
