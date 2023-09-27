[![Docs](https://readthedocs.org/projects/calamary/badge/?version=latest)](https://docs.calamary.net/en/latest/?badge=latest)

# Calamary - Forwarding- & Filtering-Proxy

Calamary is a [squid](http://www.squid-cache.org/)-like proxy.

Its focus is set on security filtering for HTTPS.

**It will not**:
* act as caching proxy
* act as reverse proxy

**Features**:
* certificate verification
* mode to enforce TLS (*deny any unencrypted connections*)
* detect plain HTTP and respond with generic HTTPS-redirect
* support for [proxy-protocol](https://github.com/pires/go-proxyproto)


## Documentation

[docs.calamary.net](https://docs.calamary.net)

## Roadmap

- [ ] Forwarding

  - [ ] TCP
  - [ ] TLS
  - [ ] HTTP
  - [ ] UDP

- [x] YAML-based configuration

- [ ] Parsing

  - [x] TCP
  - [ ] TLS
  - [ ] HTTP
  - [ ] UDP
  - [ ] DNS

- [ ] Filtering

  - [x] TCP
  - [x] TLS

    - [ ] Certificate validation

  - [x] HTTP
  - [ ] UDP

  - [ ] Matches

    - [x] Config
    - [x] Matching
    - [ ] Additional checks
