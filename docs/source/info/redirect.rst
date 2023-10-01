.. _redirect:

.. include:: ../_inc/head.rst

.. include:: ../_inc/in_progress.rst

################
Redirect Traffic
################

Basics
######

You may want/need to redirect traffic to the proxy's listeners for some use-case.

This is essential for using the :code:`transparent` mode.

For modes like :code:`proxyproto`, :code:`http`, :code:`https` or :code:`socks5` this is not necessary. (*but it's also possible using the* :ref:`Redirector <redirector>`)

You will have to choose between using **DNAT** and **TPROXY** to redirect the traffic on firewall-level.

**TProxy** has the benefit that it won't modify the packets destination. This makes processing the traffic easier and can be benefitial in regards to performance.

But it also has the drawback that traffic that originates from the proxy-server (*netfilter hook - output*) will have to be looped-back.

I personally like to use TProxy for filtering input/forward- and DNAT for output-traffic.


.. warning::
    The config-examples below may not be complete!
    If you find issues with them - please `open an issue <https://github.com/superstes/calamary/issues>`_


NFTables
########

Read more about NFTables here: `wiki.superstes.eu - NFTables <https://wiki.superstes.eu/en/latest/1/network/firewall_nftables.html>`_

DNAT
====

.. code-block:: bash

    define PROXY_PORT=4128
    define PROXY_UID=13  # user-id of the proxy-user; anti-loop
    define PROXY_PORTS={ 80, 443, 587, 25, 53, 853, 123 }

    table inet default {
        chain prerouting_dnat {
            type nat hook prerouting priority dstnat; policy accept;

            # redirect traffic from outside
            meta l4proto tcp redirect to $PROXY_PORT
            # redirect to = equivalent to 'dnat to 127.0.0.1/::1'
        }

        chain output_dnat {
            type nat hook output priority -100; policy accept;

            # redirect traffic from this host
            meta l4proto tcp meta skuid != $PROXY_UID redirect to $PROXY_PORT
        }

        chain input {
            type filter hook output priority 0; policy drop;

            # allow traffic to proxy
            meta l4proto tcp dport $PROXY_PORT ip daddr 127.0.0.1 accept comment "Allow Network to proxy"
            meta l4proto tcp dport $PROXY_PORT ip6 daddr ::1 accept comment "Allow Network to proxy"
        }

        chain output {
            type filter hook output priority 0; policy drop;

            # allow traffic to proxy
            meta l4proto tcp dport $PROXY_PORT ip daddr 127.0.0.1 accept comment "Allow localhost to proxy"
            meta l4proto tcp dport $PROXY_PORT ip6 daddr ::1 accept comment "Allow localhost to proxy"

            # optional logging
            meta l4proto tcp dport $PROXY_PORTS meta skuid $PROXY_UID ct state new log prefix "NFTables Proxy outgoing "
            # allow traffic from proxy
            meta l4proto tcp dport $PROXY_PORTS meta skuid $PROXY_UID accept comment "Allow proxy traffic"
        }
    }


TProxy
======

Full TProxy example: `gist.github.com/superstes - TProxy NFTables <https://gist.github.com/superstes/6b7ed764482e4a8a75334f269493ac2e>`_

You might need to enable some nftables kernel modules: `Kernel docs - NFTables extensions <https://docs.kernel.org/networking/tproxy.html#iptables-and-nf-tables-extensions>`_   

IPTables
########


DNAT
====

.. code-block:: bash

    PROXY_PORT=4128
    PROXY_UID=13  # user-id of the proxy-user; anti-loop

    # redirect traffic from outside
    iptables -t nat -A PREROUTING -p tcp -j REDIRECT --to-destination --to-port "$PROXY_PORT"
    # redirect traffic from localhost
    iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner "$PROXY_UID" -j REDIRECT --to-port "$PROXY_PORT"

    # allow traffic to proxy
    #   iptable -A INPUT -p tcp -d 127.0.0.1 --dport "$PROXY_PORT" -j ACCEPT
    #   iptable -A OUTPUT -p tcp -d 127.0.0.1 --dport "$PROXY_PORT" -j ACCEPT

    # optional logging
    iptables -A OUTPUT -p tcp -m conntrack --ctstate NEW -j LOG –log-prefix "IPTables Proxy outgoing "
    # allow traffic from proxy
    iptable -A OUTPUT -p tcp -m owner ! --uid-owner "$PROXY_UID" -j ACCEPT


TProxy
======

Full TProxy example: `gist.github.com/superstes - TProxy IPTables <https://gist.github.com/superstes/c4fefbf403f61812abf89165d7bc4000>`_

You might need to enable some iptables kernel modules: `Kernel docs - IPTables extensions <https://docs.kernel.org/networking/tproxy.html#iptables-and-nf-tables-extensions>`_   


TProxy
######

To run Calamary as `TPROXY <https://docs.kernel.org/networking/tproxy.html>`_ target - you will have to set `CAP_NET_RAW <https://man7.org/linux/man-pages/man7/capabilities.7.html>`_:

::

  bind to any address for transparent proxying

You can add it like this:

.. code-block:: bash

    setcap cap_net_raw=+ep /usr/bin/calamary

    # make sure only wanted users can execute the binary!
    chown root:proxy /usr/bin/calamary
    chmod 750 /usr/bin/calamary

Read more about TPROXY here:

* `wiki.superstes.eu - NFTables - TProxy <https://wiki.superstes.eu/en/latest/1/network/firewall_nftables.html#tproxy>`_

* `kernel docs <https://docs.kernel.org/networking/tproxy.html>`_

Output loopback
===============

You will have to configure a loopback route if you want to proxy 'output' traffic:

.. code-block:: bash

    echo "200 proxy_loopback" > /etc/iproute2/rt_tables.d/proxy.conf

    # These need to be configured persistent: (maybe use an interface up-hook)
    ip rule add fwmark 200 table proxy_loopback
    ip -6 rule add fwmark 200 table proxy_loopback
    ip route add local 0.0.0.0/0 dev lo table proxy_loopback
    ip -6 route add local ::/0 dev lo table proxy_loopback

    # can be checked using:
    ip rule list
    ip -6 rule list
    ip -d route show table all

    # you might need to set a sysctl:
    sysctl -w net.ipv4.conf.all.route_localnet=1

    # you might want to block 127.0.0.1 on non loopback interfaces if you enable it:
    iptables -t raw -A PREROUTING ! -i lo -d 127.0.0.0/8 -j DROP
    iptables -t mangle -A POSTROUTING ! -o lo -s 127.0.0.0/8 -j DROP
