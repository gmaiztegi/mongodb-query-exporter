# Prometheus MongoDB query exporter
![.github/workflows/action.yml](https://github.com/raffis/mongodb-query-exporter/workflows/.github/workflows/action.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/raffis/mongodb-query-exporter)](https://goreportcard.com/report/github.com/raffis/mongodb-query-exporter)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/raffis/mongodb-query-exporter?tab=subdirectories)](https://pkg.go.dev/github.com/raffis/mongodb-query-exporter?tab=subdirectories)
[![Coverage Status](https://coveralls.io/repos/github/raffis/mongodb-query-exporter/badge.svg?branch=master)](https://coveralls.io/github/raffis/mongodb-query-exporter?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/raffis/mongodb-query-exporter.svg?maxAge=604800)](https://hub.docker.com/r/raffis/mongodb-query-exporter)

MongoDB aggregation query exporter for [Prometheus](https://prometheus.io).

## Features

* Support for gauge and counter metrics
* Multiple metrics for different db/collections
* Pull and Push (Push is only supported for MongoDB >= 3.6)
* Supports multiple MongoDB servers
* Public API for Golang

Note that this is not designed to be a replacement for the [MongoDB exporter](https://github.com/percona/mongodb_exporter) to instrument MongoDB internals. This application exports custom MongoDB metrics in the prometheus format based on the queries (aggregations) you want.

## Beta notice

This software is currently beta and the API/configuration may break without notice until a stable version is released.

## Installation

Get Prometheus MongoDB aggregation query exporter, either as a [binary](https://github.com/raffis/mongodb-query-exporter/releases/latest) or packaged as a [Docker image](https://hub.docker.com/r/raffis/mongodb-query-exporter).

### Helm Chart
For kubernetes users there is an official helm chart for the MongoDB query exporter.

## Usage

```
$ mongodb_query_exporter
```

Use the `-help` flag to get help information.

If you use [MongoDB Authorization](https://docs.mongodb.org/manual/core/authorization/), best practices is to create a dedicated readonly user:

1. Create a user with '*read*' on your database, like the following (*replace username/password/db!*):

```js
db.getSiblingDB("admin").createUser({
    user: "mongodb_query_exporter",
    pwd: "secret",
    roles: [
        { role: "read", db: "mydb" }
    ]
})
```

2. Set environment variable `MONGODB_URI` before starting the exporter:

```bash
export MDBEXPORTER_MONGODB_URI=mongodb://mongodb_query_exporter:secret@localhost:27017
```

If you use [x.509 Certificates to Authenticate Clients](https://docs.mongodb.com/manual/tutorial/configure-x509-client-authentication/), pass in username and `authMechanism` via [connection options](https://docs.mongodb.com/manual/reference/connection-string/#connections-connection-options) to the MongoDB uri. Eg:

```
mongodb://CN=myName,OU=myOrgUnit,O=myOrg,L=myLocality,ST=myState,C=myCountry@localhost:27017/?authMechanism=MONGODB-X509
```

## Access metrics
The metrics are exposed at `/metrics`.

```
curl localhost:9412/metrics
```

## Configuration

The exporter is looking for a configuration in `~/config.yaml` and `/etc/mongodb-query-exporter/config.yaml` or if set the path from the env `MDBEXPORTER_CONFIG`.

You may also use env variables to configure the exporter:

| Env variable             | Description                              |
|--------------------------|------------------------------------------|
| MDBEXPORTER_CONFIG       | Custom path for the configuration        |
| MDBEXPORTER_MONGODB_URI  | The MongoDB connection URI               |
| MDBEXPORTER_MONGODB_QUERY_TIMEOUT | Timeout until a MongoDB operations gets aborted |
| MDBEXPORTER_LOG_LEVEL    | Log level                                |
| MDBEXPORTER_LOG_ENCODING | Log format                               |
| MDBEXPORTER_BIND         | Bind address for the HTTP server         |

Note if you have multiple collectors you can inject an env variable for the MongoDB connection URI like:

1. `MDBEXPORTER_COLLECTOR_0_MONGODB_URI=mongodb://srv1:27017`
2. `MDBEXPORTER_COLLECTOR_1_MONGODB_URI=mongodb://srv2:27017`
3. ...

### Format v2.0

The config format v2.0 is not supported in any version before `v1.0.0-beta5`. Please use v1.0 or upgrade to the latest version otherwise.

Example:
**`config.yml`**

```yaml
version: 2.0
bind: 0.0.0.0:9412
log:
  encoding: json
  level: info
  development: false
  disableCaller: false
collectors:
- mongodb:
    uri: mongodb://localhost:27017
    queryTimeout: 10
    maxConnection: 3
    defaultInterval: 5
  metrics:
  - name: myapp_example_simplevalue_total
    type: gauge
    help: 'Simple gauge metric'
    value: total
    labels: []
    mode: pull
    interval: 10
    constLabels: []
    database: mydb
    collection: objects  
    pipeline: |
      [
        {"$count":"total"}
      ]  
  - name: myapp_example_processes_total
    type: gauge
    help: 'The total number of processes in a job queue'
    value: total
    mode: push
    labels: [type,status]
    constLabels: []
    database: mydb
    collection: queue
    pipeline: |
      [
        {"$group": {
          "_id":{"status":"$status","name":"$class"},
          "total":{"$sum":1}
        }},
        {"$project":{
          "_id":0,
          "type":"$_id.name",
          "total":"$total",
          "status": {
            "$switch": {
                "branches": [
                   { "case": { "$eq": ["$_id.status", 0] }, "then": "waiting" },
                   { "case": { "$eq": ["$_id.status", 1] }, "then": "postponed" },
                   { "case": { "$eq": ["$_id.status", 2] }, "then": "processing" },
                   { "case": { "$eq": ["$_id.status", 3] }, "then": "done" },
                   { "case": { "$eq": ["$_id.status", 4] }, "then": "failed" },
                   { "case": { "$eq": ["$_id.status", 5] }, "then": "canceled" },
                   { "case": { "$eq": ["$_id.status", 6] }, "then": "timeout" }
                ],
                "default": "unknown"
            }}
        }}
      ]
```

### Format v1.0

The config version v1.0 is the predescer of v2.0 and does not have support for multiple MongoDB servers
nor is it possible to customize logging.
When possible use v2.0 however v1.0 support won't be dropped.

Example:
**`config.yml`**

```yaml
version: 1.0
bind: 0.0.0.0:9412
logLevel: info
mongodb:
  uri: mongodb://localhost:27017
  connectionTimeout: 3
  maxConnection: 3
  defaultInterval: 5
metrics:
- name: myapp_example_simplevalue_total
  type: gauge
  help: 'Simple gauge metric'
  value: total
  labels: []
  mode: pull
  interval: 10
  constLabels: []
  database: mydb
  collection: objects  
  pipeline: |
    [
      {"$count":"total"}
    ]  
- name: myapp_example_processes_total
  type: gauge
  help: 'The total number of processes in a job queue'
  value: total
  mode: push
  labels: [type,status]
  constLabels: []
  database: mydb
  collection: queue
  pipeline: |
    [
      {"$group": {
        "_id":{"status":"$status","name":"$class"},
        "total":{"$sum":1}
      }},
      {"$project":{
        "_id":0,
        "type":"$_id.name",
        "total":"$total",
        "status": {
          "$switch": {
              "branches": [
                 { "case": { "$eq": ["$_id.status", 0] }, "then": "waiting" },
                 { "case": { "$eq": ["$_id.status", 1] }, "then": "postponed" },
                 { "case": { "$eq": ["$_id.status", 2] }, "then": "processing" },
                 { "case": { "$eq": ["$_id.status", 3] }, "then": "done" },
                 { "case": { "$eq": ["$_id.status", 4] }, "then": "failed" },
                 { "case": { "$eq": ["$_id.status", 5] }, "then": "canceled" },
                 { "case": { "$eq": ["$_id.status", 6] }, "then": "timeout" }
              ],
              "default": "unknown"
          }}
      }}
    ]
```

## Go API
Instead using the mongodb-query-exporter you may use the API to integrate the exporter within your go project.
Please check out the [go package reference](https://pkg.go.dev/badge/github.com/raffis/mongodb-query-exporter?tab=subdirectories).

## Used by
* The balloon helm chart implements the mongodb-query-exporter to expose general stats from the MongoDB like the number of total nodes or files stored internally or externally.
See the [config-map here](https://github.com/gyselroth/balloon-helm/blob/master/unstable/balloon/charts/balloon-mongodb-metrics/templates/config-map.yaml).


Please submit a PR if your project should be listed here!
