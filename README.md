# Trewoga v0.0.2

Really tiny and simple service monitoring system

## Overview

Trewoga runs a tiny server that listens for pings from monitored services. If
there is no pings from one of the services for specified interval of time, it is
considered to be failed and all subscribed users are notified via Telegram.

## Installation

1. Run `make` to compile the project
2. Copy `bin/amd64/trewoga` binary to your server
3. Run `trewoga --config /path/to/trewoga.yml` on your server

## Configuration

(see sample configuration file `trewoga.sample.yml`)

* `address` (required) - URL to be used for Trewoga API. If HTTPS address is
used, certificates are configured automatically (using
[Let's Encrypt](https://letsencrypt.org)).
* `db` (required) - path to the sqlite database
* `token` (required) - Telegram Bot API token
* `proxy` (optional) - HTTPS proxy for Telegram (useful in countries where
Telegram is blocked)
* `cert` (options, default = ".") - path for Let's Encrypt to store certificates

## Managing services

Service management is done using the CLI utility:

* list of all configured services
```
# trewoga --config /path/to/trewoga.yml service:list
name                   =  service1
token                  =  eaf49b34-a254-4127-afb5-7ce25efa4c1c
failure                =  true
recovering             =  false
maintenance            =  false
failure timeout        =  10s
recovery interval      =  2m0s
maintenance timeout    =  20m0s
ping at                =  0001-01-01 00:00:00 +0000 UTC
recovery at            =  0001-01-01 00:00:00 +0000 UTC
maintenance failure at =  0001-01-01 00:00:00 +0000 UTC
```

* registering a new service, following parameters can be specified:
  * `failure` (seconds, default 30) - if a service delays the ping for more than
  failure timeout, it is considered to be failed (subscribed users are notified
  via Telegram)
  * `recovery` (seconds, default 120) - if a failed service sends pings without
  failures for more than recovery interval, it is considered not to be failed
  anymore (subscribed users are notified via Telegram)
  * `maintenance` (minutes, default 20) - service can be put to maintenance mode
  via API. When in maintenance mode, services are not checked for failures.
  However if the service is in maintenance longer than maintenance timeout, it
  is considered to be failed (subscribed users are notified via Telegram)
```
# trewoga --config /path/to/trewoga.yml service:add service1 --failure 10 --recovery 120 --maintenance 20
name                   =  service1
token                  =  eaf49b34-a254-4127-afb5-7ce25efa4c1c
failure                =  true
recovering             =  false
maintenance            =  false
failure timeout        =  10s
recovery interval      =  2m0s
maintenance timeout    =  20m0s
ping at                =  0001-01-01 00:00:00 +0000 UTC
recovery at            =  0001-01-01 00:00:00 +0000 UTC
maintenance failure at =  0001-01-01 00:00:00 +0000 UTC
```

## API

* ping
```
curl -X POST https://trewoga.example.com/v1/services/SERVICE_TOKEN/ping
```

* enable maintenance mode
```
curl -X POST https://trewoga.example.com/v1/services/SERVICE_TOKEN/maintenance
```

* disable maintenance mode
```
curl -X DELETE https://trewoga.example.com/v1/services/SERVICE_TOKEN/maintenance
```
