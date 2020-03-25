# MailProxy

A dead simple Mail forward proxy written in Golang.

## HowTo guide

### Docker

- Build it

```bash
$ docker build -t go-mailproxy:v1.0 .
```

- Generate cert key files, then put it in `/path/to/directory/contains/cert` (Change this path with the actual path). For example:

```bash
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
# https://safecurves.cr.yp.to/
# List ECDSA the supported curves (openssl ecparam -list_curves)
$ openssl ecparam -genkey -name secp384r1 -out server.key
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

- Update `config.yml`, you could refer [an example config](./etc/config.yml)

- Run it (build command is suppored as well, check [Makefile](./Makefile) for more details)

```bash
$ make run
```

### Manual

- Generate cert key files, then put it in `/path/to/directory/contains/cert` (Change this path with the actual path). For example:

```bash
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
# https://safecurves.cr.yp.to/
# List ECDSA the supported curves (openssl ecparam -list_curves)
$ openssl ecparam -genkey -name secp384r1 -out server.key
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

- Update config file.

- Run it.

```bash
$ cd mailproxy/
$ go build -mod vendor -o bin/mailproxy
$ ./bin/mailproxy -conf /path/to/directory/contains/config
```
