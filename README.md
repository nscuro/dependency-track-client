# dependency-track-client

![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)

## Installation

`GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...`

## Usage

```
$ ./dtrack -h
Usage:
  dtrack [command]

Available Commands:
  audit       Audit for vulnerabilities
  help        Help about any command
  report      Generate a vulnerability report

Flags:
  -k, --api-key string    dependency-track api key
  -u, --base-url string   dependency-track base url
  -h, --help              help for dtrack

Use "dtrack [command] --help" for more information about a command.
```

Global flags can be provided via environment variables as well:

```
$ export DTRACK_BASE_URL=https://dependencytrack.evilcorp.com
$ export DTRACK_API_KEY=0sl67mjen99zxb2y
```

### Audit

```
$ ./dtrack audit --project Dependency-Track --version 3.8.0 --bom ./bom.xml --autocreate
```

### Report

```
$ ./dtrack report --project Dependency-Track --version 3.8.0 --template ./examples/report.tpl --output report.html
```
