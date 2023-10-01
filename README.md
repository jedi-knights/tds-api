# TopDrawerSoccer API

![GitHub Actions](https://github.com/jedi-knights/tds-api/workflows/CI/badge.svg)


## Setup

### echo-swagger

In order to get started we need to first install the command line tool.

```bash
go get -u github.com/swaggo/swag/cmd/swag
```

Then I needed to initialize the root like so

```bash
swag init -g cmd/server/main.go
```

This generate the following files in the docs directory:

* docs/docs.go
* docs/swagger.json
* docs/swagger.yaml

Before I did this I was getting a 500 error while the swagger UI tried to access doc.json.

I've created this repo with the following directories under the pkg directory.

* handlers (contains the http handlers)
* models (contains the data models)
* services (contains the business logic)

## Docker

Running the image

> docker run -i -t -p 3000:8080 tds-api

This will map to port 3000 on the host machine.

The options -i and -t allows us to execute the container into the interactive mode that allows us to shut it down with CTRL+C when needed.

The option -p 8080:8080 maps the container port 8080 to the same port on our machine. It will effectively allow us to talk with the echo server using the port 8080.

### Pruning Docker Images

Because docker images can take up a lot of space on your machine, it's a good idea to prune them 
from time to time.  It will free up space on your machine.

```bash
docker image prune -a -f
```

or by using the Makefile

```bash
make docker-prune
```

## Cross Platform Builds

To get a listing of the possible platforms that we can build for, we can run the following command:

```bash
go tool dist list
```

## References

* [echo-swagger](https://github.com/swaggo/echo-swagger)
* [Optimizing Golang Docker images with multi-stage build](https://medium.com/geekculture/optimizing-golang-docker-images-with-multi-stage-builds-ca7b305faa)
* 