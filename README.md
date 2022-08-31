# ifconfig - Public IP address API

[![cicd](https://github.com/jakewmeyer/ifconfig/workflows/CICD/badge.svg)](https://github.com/jakewmeyer/ifconfig/actions?query=workflow%3ACICD)

## Usage

### Plaintext

```bash
curl -s 'https://ifconfig.jakemeyer.sh'
```

```text
192.168.1.1
```

### JSON

```bash
curl -s 'https://ifconfig.jakemeyer.sh/json' | jq
```

```json
{"ip":"192.168.1.1"}
```

## Setup

### Docker Compose v3+

`docker-compose.yml`

```yaml
services:
  ifconfig:
    container_name: "ifconfig"
    image: ghcr.io/jakewmeyer/ifconfig:latest
    ports:
      - "7000:7000"
    restart: "unless-stopped"
```

```bash
docker-compose up
```

### Docker run

```bash
docker run -p 7000:7000 ghcr.io/jakewmeyer/ifconfig:latest
```
