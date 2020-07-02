# dependency-track-client

![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)

## Installation

`GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...`

## Usage

```
$ export DTRACK_BASE_URL=https://dependencytrack.evilcorp.com
$ export DTRACK_API_KEY=...
```

### Audit

`$ ./dtrack audit --project Dependency-Track --version 3.8.0 --bom ./bom.xml --autocreate`

### Report

`$ ./dtrack report --project Dependency-Track --version 3.8.0 --template ./examples/report.tpl --output report.html`
