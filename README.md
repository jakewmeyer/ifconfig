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
