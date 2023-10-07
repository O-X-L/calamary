[![Status Lint](https://github.com/superstes/calamary/actions/workflows/lint.yml/badge.svg?branch=latest)](https://github.com/superstes/calamary/actions/workflows/lint.yml)
[![Status Unit Tests](https://github.com/superstes/calamary/actions/workflows/test.yml/badge.svg?branch=latest)](https://github.com/superstes/calamary/actions/workflows/test.yml)
[![Status Unit Test Coverage](https://codecov.io/gh/superstes/calamary/graph/badge.svg?token=PPNLDDS0M8)](https://codecov.io/gh/superstes/calamary)
[![Status Integration Tests](https://badges.calamary.net/latest.calamary.test.svg)](https://github.com/superstes/calamary/tree/latest/test)
[![Status Documentation](https://readthedocs.org/projects/calamary/badge/?version=latest)](https://docs.calamary.net/en/latest/)

# Calamary - Forwarding- & Filtering-Proxy

Calamary is a [squid](http://www.squid-cache.org/)-like proxy.

Its focus is set on transparent security filtering.

## Contributing

Feel free to contribute to this project!

[Reporting issues](https://github.com/superstes/calamary/issues), [discussing implementation](https://github.com/superstes/calamary/discussions), [extending documentation](https://github.com/superstes/calamary/tree/latest/docs) and adding unit-/[integration-tests](https://github.com/superstes/calamary/tree/latest/test) is very welcome!

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

  - [ ] Authentication

- [ ] Forwarding

  - [x] TCP

    - [x] HTTP

  - [x] TLS

    - [ ] TLS Interception

  - [ ] UDP

    - [ ] QUIC

- [x] YAML-based configuration

- [ ] Parsing

  - [ ] Basic

    - [x] TCP
    - [ ] UDP

      - [ ] QUIC

    - [x] TLS

      - [ ] [ECH](https://datatracker.ietf.org/doc/draft-ietf-tls-esni/)/ESNI handling (*encrypted SNI*)

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
