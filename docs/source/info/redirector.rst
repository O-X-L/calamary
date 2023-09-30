.. _redirector:

.. |proxy_forwarder| image:: https://wiki.superstes.eu/en/latest/_images/squid_remote.png
   :class: wiki-img

.. include:: ../_inc/head.rst

##########
Redirector
##########

The redirector will be a smaller version of Calamary. (*without filtering*)

It can be used to forward traffic from remote locations to your proxy servers.

|proxy_forwarder|

Per example - this might be useful if you have:

* Distributed systems

  * Cloud servers that are only connected to public WAN and should send all their outbound traffic over your proxy

  * Simple/dumb firewalls/routers that should send the networks outbound traffic over your proxy

As it utilizes the commonly used `proxy-protocol <https://www.haproxy.com/blog/use-the-proxy-protocol-to-preserve-a-clients-ip-address>`_ to redirect the traffic, it might also be useful in combination with other proxies.
