.. _intro:

.. |proxy_forwarder| image:: https://wiki.superstes.eu/en/latest/_images/squid_remote.png
   :class: wiki-img

.. include:: ../_inc/head.rst

#####
Intro
#####

Calamary is a `squid <http://www.squid-cache.org/>`_-like proxy.

Its focus is set on **security filtering for HTTPS**.

**It will not**:

* act as caching proxy
* act as reverse proxy

**Features**:

* basic traffic filtering - see :ref:`Rules <rules>`
* certificate verification
* enforce TLS (*deny any unencrypted connections*)
* detect plain HTTP and respond with generic HTTPS-redirect
* intercept-, http-proxy and proxy-proto-modes

  * support for `proxy-protocol <https://github.com/pires/go-proxyproto>`_

Getting Started
###############

 :ref:`Getting started <getting_started>`

Why?
####

Forward proxies are very useful to enforce a security-baseline in networks and a must-have for Zero-Trust environments.

Many enterprises and individuals will use proxies integrated with vendor network-firewalls or cloud-services to handle this filtering.

But some of us might like to keep control over that system.

The usage of go-based applications is easy (*single binary*) and can perform well.

Why not use Squid?
==================

**Squid has some limitations** that make its usage more complicated than it should be.

**Per example**:

* intercept/transparent mode - no native solution for `the DNAT restrictions <http://www.squid-cache.org/Advisories/SQUID-2011_1.txt>`_

  Related errors:

  * `NF getsockopt(ORIGINAL_DST) failed`
  * `NAT/TPROXY lookup failed to locate original IPs`
  * `Forwarding loop detected`


* intercept/transparent mode - `host verification - using DNS <http://www.squid-cache.org/Doc/config/host_verify_strict/>`_

  does hit issues with todays DNS-handling of major providers:

  * TTLs <= 1 min (*p.e. download.docker.com, debian.map.fastlydns.net*)

  Related error: `Host header forgery detected`


Squid is a good and stable software. But I get the feeling it needed to grow into more than it was designed for initially. Some behavior is inconsistent between modes and not optimized for todays IT-world.

I would much preferr a keep-it-simple approach. Even if that means that some nice-to-have features are not implemented.


How?
####

* Use TLS-SNI as target instead of HTTP Host-Header


* Optionally use additional DNS-based verfication if TTL > 3 min


* Whenever it is not possible to route the traffic through the proxy..

  To overcome the DNAT restriction, of losing the real target IP, the proxy will have a lightweight forwarder mode:

  |proxy_forwarder|


* **Transparent traffic interception will be the focus**.

  Setting the environment-variables 'HTTP_PROXY', 'HTTPS_PROXY', 'http_proxy' and 'https_proxy' for all applications and HTTP-clients may be problematic/too inconsistent
