# Calamary - Forwarding- & Filtering-Proxy

Calamary is a [squid](http://www.squid-cache.org/)-like proxy.

Its focus is set on security filtering for HTTPS.

**It will not**:
* act as caching proxy
* act as reverse proxy

**Features**:
* certificate verification
* detect plain HTTP and respond with generic HTTPS-redirect


## Why?

**Squid has some limitations** that make its usage more complicated than it should be.

**Per example**:

* intercept/transparent mode - [DNAT not supported](http://www.squid-cache.org/Advisories/SQUID-2011_1.txt)

  Related errors:

  * `NF getsockopt(ORIGINAL_DST) failed`
  * `NAT/TPROXY lookup failed to locate original IPs`
  * `Forwarding loop detected`


* intercept/transparent mode - [host verification - using DNS](http://www.squid-cache.org/Doc/config/host_verify_strict/)

  does hit issues with todays DNS-handling of major providers:

  * TTLs <= 1 min (*p.e. download.docker.com, debian.map.fastlydns.net*)

  Related error: `Host header forgery detected`



## How?

* Use TLS-SNI as target instead of HTTP Host-Header


* Optionally use additional DNS-based verfication if TTL > 3 min


* Whenever it is not possible to route the traffic through the proxy..

  To overcome the DNAT restriction, of losing the real target IP, the proxy will have a lightweight forwarder mode:

  <img src="https://wiki.superstes.eu/en/latest/_images/squid_remote.png" alt="Proxy forwarder" width="400">


* Transparent traffic interception will be the focus.

  Setting the environment-variables 'HTTP_PROXY', 'HTTPS_PROXY', 'http_proxy' and 'https_proxy' for all applications and HTTP-clients may be problematic/too inconsistent


## Roadmap

- [ ] Forwarding

  - [ ] TCP
  - [ ] TLS
  - [ ] HTTP
  - [ ] UDP

- [ ] YAML-based configuration

- [ ] Filtering

  - [ ] TCP
  - [ ] TLS
  - [ ] HTTP
  - [ ] UDP

  - [ ] Certificate validation
  - [ ] ACLs

    - [ ] Matching
    - [ ] Additional checks
