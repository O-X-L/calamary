.. _rules:

.. include:: ../_inc/head.rst

.. include:: ../_inc/in_progress.rst

#####
Rules
#####

You can define a list of rules that calamary will apply to the traffic passing through it.

Rules basically consist of a :code:`match` and an :code:`action`!

.. code-block:: yaml

    rules:
      - match:
          dest: '192.168.100.0/24'
        action: 'drop'

Matches
#######

Multiple matches can be defined in a single rule.

.. code-block:: yaml

    - match:
        src: 'IP OR NET+CIDR'
        dest: 'IP OR NET+CIDR'
        port: 'NUMBER'  # destination ports
        sport: 'NUMBER'  # source ports
        protoL3: 'ip4/ivp4/ip6/ip6'
        protoL4: 'tcp/udp'  # others might be supported later on
        protoL5: 'tls/http/dns/ntp'  # others might be supported later on
        dns: 'DOMAIN'  # domain/TLS-SNI to match
        encrypted: 'true/false/yes/no'  # match TLS traffic

The value of matches is **case-insensitive** by default.

NOTE: The HTTP host-header domain is not compared if :code:`dns` is used - as it can be modified easily.

You can define **multiple values** for each match.

Matches can also be **negated** by using the :code:`!` prefix:

.. code-block:: yaml

    rules:
      - match:
          port: ['!80', '!443', '!587']
        action: 'drop'

      - match:
          dest: '!192.168.0.0/16'
          port: 443
          protoL4: 'tcp'
        action: 'accept'

Packets that don't match any :code:`accept` rule will be **dropped by default**.

Actions
#######

Available actions include:

* 'accept' (*alias: 'allow'*)
* 'deny' (*alias: 'drop'*)

Other actions like 'limit' will be implemented later on.

Variables
#########

Calamary enables you to define variables that can be used inside your ruleset.

.. code-block:: yaml

    vars:
      - name: 'net_private'
        value: ['192.168.0.0/16', '172.16.0.0/12', '10.0.0.0/8']
      - name: 'svc_http'
        value: [80, 443]

Variables are referenced using the :code:`$` prefix.

Whenever you use a variable, you can also negate it like any other value:

.. code-block:: yaml

    rules:
      - match:
          src: '$net_private'
          dest: '!$net_private'
        action: 'accept'
