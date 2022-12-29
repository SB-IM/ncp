# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.3.0] - 2022-12-29

### Add

* Add `os.env` expand

```yaml
mqttd:
  id: ${NCP_UUID}
  static:
    type: ${NCP_TYPE}
  broker: ${NCP_MQTT_URL}
```

## [2.2.0] - 2022-04-01

### Add

* Add `type` in config
* Add goreleaser

### Change

* Change (config) `Link_id` int to string

### Fix

* Fix mqtt auth

## [2.1.3] - 2021-07-12

### Add

* Add docker container

### Change

* Improve build
* Print log

### Fix

* Linux build depend 'libc'

## [2.1.2] - 2021-06-12

* Add License MPL v2
* Open source first release

## [2.1.1] - 2021-04-30

### Fixed

* mqtt onServerDisconnect Panic

## [2.1.0] - 2021-04-24

### Added

* [README_zh.md](README_zh.md)
* ncpio Debug mode stdout
* mqtt rpc LRUCache
* mqtt rpc config QoS
* conf/template

### Changed

* Disable restore session

### Fixed

* mqtt jsonrpc Notify

## [2.0.0] - 2021-01-26

### Added

* First release

