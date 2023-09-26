.. _getting_started:

.. include:: _inc/head.rst

.. include:: _inc/in_progress.rst


###############
Getting Started
###############

Installation
############

Download the binary for your desired target system from the `latest relase page <https://github.com/superstes/calamary/releases/latest>`_

Run
###

Just execute it.

.. code-block:: bash

    /usr/bin/calamary


Config-validation only:

.. code-block:: bash

    /usr/bin/calamary -v


TProxy
======

To run Calamary as `TPROXY <https://docs.kernel.org/networking/tproxy.html>`_ target - you will have to set [CAP_NET_RAW](https://man7.org/linux/man-pages/man7/capabilities.7.html):

> bind to any address for transparent proxying

You can add it like this:

```bash
setcap cap_net_raw=+ep /usr/bin/calamary

# make sure only wanted users can execute the binary!
chown root:proxy /usr/bin/calamary
chmod 750 /usr/bin/calamary
```

Read more about TPROXY here: `wiki.superstes.eu - NFTables - TProxy <https://wiki.superstes.eu/en/latest/1/network/nftables.html#tproxy>`_

Systemd service
===============

Example systemd service to run Calamary:


.. code-block:: text

    [Unit]
    Description=Calamary Forward Proxy
    Documentation=https://docs.calamary.net
    Documentation=https://github.com/superstes/calamary
    After=network-online.target
    Wants=network-online.target

    [Service]
    Type=simple
    Environment=CONFIG=/etc/calamary/config.yml

    # validate before start/restart
    ExecStartPre=/usr/bin/calamary -f $CONFIG -v
    ExecStart=/usr/bin/calamary -f $CONFIG

    # validate before reload
    ExecReload=/usr/bin/calamary -f $CONFIG -v
    ExecReload=/bin/kill -HUP $MAINPID

    User=proxy
    Group=proxy
    Restart=on-failure
    RestartSec=5s

    StandardOutput=journal
    StandardError=journal
    SyslogIdentifier=calamary

    [Install]
    WantedBy=multi-user.target


Configuration
#############

See :ref:`Rules <rules>`for more details about defining the filter ruleset.

The default config path is :code:`/etc/calamary/config.yml`

Basic config example:

.. code-block:: yaml

    ---

    service:
    listen:
      port: 4128
      ip4:
        - '127.0.0.1'
      ip6:
        - '::1'
      tcp: true
      udp: false  # not yet implemented
      transparent: false  # tproxy mode

    debug: false
    timeout:
      connection: 5
      handshake: 5
      dial: 5
      intercept: 2

    output:
      fwmark: 0
      interface: ''

    vars:
      - name: 'net_private'
        value: ['192.168.0.0/16', '172.16.0.0/12', '10.0.0.0/8']
      - name: 'svc_http'
        value: [80, 443]

    rules:
      - match:
          dest: '192.168.100.0/24'
        action: 'drop'

      - match:
          port: ['!443', '!80']
        action: 'drop'

      - match:
          src: '$net_private'
          dest: '$net_private'
          port: '$svc_http'
          protoL4: 'tcp'
        action: 'accept'

      - match:
          dest: '!$net_private'
          port: 443
          protoL4: 'tcp'
        action: 'accept' 
 