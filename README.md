[![Unit Tests](https://github.com/superstes/calamary/actions/workflows/test.yml/badge.svg?branch=latest)](https://github.com/superstes/calamary/actions/workflows/test.yml)
[![Test Coverage](https://codecov.io/gh/superstes/calamary/graph/badge.svg?token=UD4TM8N256)](https://codecov.io/gh/superstes/calamary)
[![Lint](https://github.com/superstes/calamary/actions/workflows/lint.yml/badge.svg?branch=latest)](https://github.com/superstes/calamary/actions/workflows/lint.yml)
[![Documentation](https://readthedocs.org/projects/calamary/badge/?version=latest)](https://docs.calamary.net/en/latest/?badge=latest)

# Calamary - Forwarding- & Filtering-Proxy

Calamary is a [squid](http://www.squid-cache.org/)-like proxy.

Its focus is set on transparent security filtering.

## Documentation

[docs.calamary.net](https://docs.calamary.net)

## Roadmap

- [ ] Listeners

  - [ ] Transparent

    - [x] TCP
    - [ ] UDP

  - [ ] Proxy-Protocol

  - [x] HTTP Proxy

  - [ ] HTTPS Proxy

  - [ ] SOCKS5 Proxy

- [ ] Forwarding

  - [x] TCP

    - [x] HTTP

  - [x] TLS

    - [ ] TLS Interception

  - [ ] UDP

- [x] YAML-based configuration

- [ ] Parsing

  - [ ] Basic

    - [x] TCP
    - [ ] UDP
    - [x] TLS
    - [ ] Identify common protocols

  - [ ] Listener-Specific

    - [ ] Proxy-Protocol
    - [x] HTTP Proxy
    - [ ] SOCKS5 Proxy

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

  - [ ] DNS
