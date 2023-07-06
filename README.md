# Reverse Proxy
## Description
This is a simple reverse proxy that could block request based in a list of blockers.
Also this reverse proxy have the capability to mask sensitive information in response body.
It only Mask the response of GET requests.

## Features

* Set blocker list in a configuration file.
* Blocker list includes: MethodBlocker, PathBlocker, ParamBlocker and HeaderBlocker.
* Includes two maskers: CreditCardMasker, EmailMasker.
* Easy to extend with new blockers and maskers.
* Log all incoming request and response in human readable format.
* Simple, only use standard library besides a logger.
* Graceful shutdown.

## How to Set Up
#### Prerequisites
* [Golang](https://golang.org/doc/install)

### Configuration
The app uses a configuration file in TOML format.
The configuration path is set via flags when running the app.
If non is provided the app will use an example configuration file located in internal/config/example_config.toml

Toml example:
```
TargetURL = "http://localhost:8080"
ReverseProxyPort = 8081
[HeaderBlocker]
  [HeaderBlocker.HeaderMap]
    X-Blocker = "Block"
    Authorization = "Secret"
[ParamBlocker]
  [ParamBlocker.ParamsMap]
    apikey = "token"
[PathBlocker]
  path = ["/admin", "/private"]
[MethodBlocker]
  method = ["POST", "PUT"]
```

#### Build the app
```
make build
```
And then run the app with the configuration file path
```
./reverse-proxy -config /folder/config.toml
```

#### Run the tests
```
make tests
```
#### Run the examples
```
make example
```

After running make example, you can test the reverse proxy with the target_server_example located in test/ folder.
```
curl -X GET -v http://localhost:8081/get 
```

## Future Work - Nice To have
* In order to be production ready it needs more work with:
  * Streaming.
  * Websockets.
  * Compresses data.
  * More testing with the mask to avoid leaks.
* Healthcheck, readiness endpoint.
* TLS support.
* Add Timeout configuration to the reverse proxy.
* Metrics.
* Benchmarking and performance testing:
  * Improve the mask regexes with others [lib](https://pkg.go.dev/github.com/flier/gohs/hyperscan) [lib1] (https://github.com/google/re2)
  * Check blockers concurrently.
* Add more testcases of strange and edgy cases.

## Contributors
* [sgrodriguez](https://github.com/sgrodriguez)
