# dependency-track-client

![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)

## Installation

`GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...`

## Usage

```
age:
  dtrack [command]

Available Commands:
  audit       Audit for vulnerabilities
  bom         Export or upload BOMs
  help        Help about any command
  report      Generate reports
  version     Display version information

Flags:
  -k, --api-key string           Dependency-Track API key
  -h, --help                     help for dtrack
      --project string           Project UUID
      --project-name string      Project name
      --project-version string   Project version
  -u, --url string               Dependency-Track URL

Use "dtrack [command] --help" for more information about a command.
```

Dependency-Track's URL and the API key can be provided via environment variables as well:

```
$ export DTRACK_BASE_URL=https://dependencytrack.evilcorp.com
$ export DTRACK_API_KEY=0sl67mjen99zxb2y
```

### Audit

```
$ ./dtrack audit \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom ./bom.xml --autocreate
```

### BOM

```
$ ./dtrack bom get \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    -o bom.xml
```

```
$ ./dtrack bom upload \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom bom.xml --autocreate
```

### Report

```
$ ./dtrack report \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --template ./examples/project-report.html \
    --output report.html
```
