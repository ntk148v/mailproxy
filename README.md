# MailProxy

A dead simple Mail forward proxy written in Golang.

## HowTo guide

### Docker

- Build it

```bash
$ docker build -t go-mailproxy:v1.0 .
```

- Update `config.yml`, you could refer [an example config](./etc/config.yml)

- Run it (build command is suppored as well, check [Makefile](./Makefile) for more details)

```bash
$ make run
```

### Manual

- Update config file.

- Run it.

```bash
$ cd mailproxy/
$ go build -mod vendor -o bin/mailproxy
$ ./bin/mailproxy -conf /path/to/directory/contains/config
```
