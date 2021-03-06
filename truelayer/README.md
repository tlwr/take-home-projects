# truelayer-take-home-pokemon-api

![ci](https://github.com/tlwr/take-home-projects/truelayer/workflows/ci/badge.svg)

An API which proxies two other APIs:

* [Pokemon API](https://pokeapi.co)
* [Shakespeare Translations](https://funtranslations.com/api/shakespeare)

The API requested API was specified
[in this Google Doc](https://docs.google.com/document/d/1OEa191OL9QF96JDkIZHWUVuiWsDMwVT810rz6SUc-dY)

## Usage

```
go build
./truelayer-take-home-pokemon-api
```

By default the API will:

* use LIVE APIs
* listen on all interfaces (do you trust my code? do you trust your network?)
* listen on port 5000

```
# get help
./truelayer-take-home-pokemon-api -help
```

```
# use stub APIs
./truelayer-take-home-pokemon-api -stubs
```

```
./truelayer-take-home-pokemon-api -port 8080 -bind 127.0.0.1
```

## Design

There is an API client for the Shakespeare Translation API

There is a `PokemonClient` interface which represents the Pokemon API.

This interface is implemented multiple times:
* for real usage: `pokemonClient struct`
* the testing: `FakePokemonClient struct` (`counterfeiter`)
* for caching
* for translating descriptions

One `PokemonClient` wraps another,
the actual server uses `caching(translating(normal client))`,
so that Pokemon are only requested once and therefore translations are requested once.

Caching is implemented because the Shakespeare API allegedly has a harsh rate limit.

## Development

Check the `Makefile`:

* `make` or `make test` to run tests
* `make generate` to generate fakes

## Docker

Theoretically this is packaged with Docker. This is untested because Docker support on M1 Mac's is non-existent.

## How to improve this

* Add GoDoc for the packages

* Expose any packages that would be useful outside of `internal`

* Add logging (in practice I would use `logrus` and a real HTTP server that supported middleware)

* Add metrics (Prometheus with a real HTTP server with middleware)
  * Request/response metrics
  * Pokemon/Shakespeare API client metrics

* Do some sanitisation on the descriptions that come from the APIs, right now any apostrophes or quotes get garbled
