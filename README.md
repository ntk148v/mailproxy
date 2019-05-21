# MailProxy

A dead simple Mail forward proxy written in Golang.

## HowTo guide

### Docker

* Build it

```bash
$ docker build -t go-mailproxy:v1.0 .
```

* Update `config.yml`, you could refer [an example config](./etc/config.yml)

* Run it

```bash
$ docker run -d -p 9011:9011 --name mailproxy -v /path/to/config.yml:/etc/mailproxy/config.yml go-mailproxy:v1.0
```

### Manual

* Update config file.

* Run it.

```bash
$ ./bin/mailproxy -conf /path/to/config.yml
```
