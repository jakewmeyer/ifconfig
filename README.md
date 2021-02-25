# ifconfig

[![cicd](https://github.com/jakewmeyer/ifconfig/workflows/CICD/badge.svg)](https://github.com/jakewmeyer/ifconfig/actions?query=workflow%3ACICD)

[![Coverage Status](https://coveralls.io/repos/github/jakewmeyer/ifconfig/badge.svg?branch=master)](https://coveralls.io/github/jakewmeyer/ifconfig?branch=master)

## Usage

### Plaintext

```http
GET https://ifconfig.jakemeyer.sh
```

```text
192.168.1.1
```

### JSON

```http
GET https://ifconfig.jakemeyer.sh?json
```

```json
{"ip":"192.168.1.1"}
```

## Setup

### Docker Compose v3+

```yaml
version: "3"

services:
  ifconfig:
    container_name: "ifconfig"
    image: ghcr.io/jakewmeyer/ifconfig:latest
    ports:
      - "7000:7000"
    restart: "unless-stopped"
```

### Docker run

```bash
docker run -p 7000:7000 ghcr.io/jakewmeyer/ifconfig:latest
```

### From source - requires go install

```bash
git clone https://github.com/jakewmeyer/ifconfig.git && cd ifconfig
```

```bash
go run main.go
```
