#!/bin/bash

TEST_TARGET='135.181.170.219'

sudo iptables -t nat -I OUTPUT -d "$TEST_TARGET" -p tcp -j DNAT --to-destination 127.0.0.1:4128
