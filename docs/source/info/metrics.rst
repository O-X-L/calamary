.. _metrics:

.. include:: ../_inc/head.rst

#######
Metrics
#######

Calamary provides a built-in `Prometheus exporter <https://prometheus.io/>`_.

You can enable it using the config-file:


.. code-block:: yaml

    ---

    service:
      ...

      metrics:
        enabled: false
        port: 9512

    ...

Test
####

.. code-block:: yaml

    curl http://localhost:9512/metrics

Metric items
############


.. code-block:: text

    # traffic
    calamary_bytes_rcv 12707
    calamary_bytes_sent 1760
    calamary_current_connections 1
    calamary_req_protoL3{protocol="IPv4"} 3
    calamary_req_protoL5{protocol="HTTP"} 1
    calamary_req_protoL5{protocol="TLS"} 2
    calamary_req_tcp 3
    calamary_req_tls_version{version="1.2"} 2
    calamary_req_tls_version{version="none"} 1

    # rules
    calamary_rule_actions{action="accept"} 3
    calamary_rule_actions{action="deny"} 1
    calamary_rule_hits{ruleId="0"} 3
    calamary_rule_hits{ruleId="1"} 3
    calamary_rule_hits{ruleId="2"} 3
    calamary_rule_hits{ruleId="3"} 3
    calamary_rule_matches{ruleId="3"} 3
