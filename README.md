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


## References

* [echo-swagger](https://github.com/swaggo/echo-swagger)