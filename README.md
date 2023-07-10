# Reverse Proxy
## Description
In addition to fulfilling the typical functionality of a reverse proxy, this tool is capable of blocking requests based on a specified list of blockers. Additionally, it has the functionality to conceal sensitive information in the response body. Notably, this feature is applied exclusively to the response of GET requests.

## Features

* Configuration file-based setting of blockers.
* Blocker list includes MethodBlocker, PathBlocker, ParamBlocker, and HeaderBlocker.
* Includes two maskers: CreditCardMasker, EmailMasker.
* Easy to extend with new blockers and maskers.
* Log all incoming requests and responses in human-readable format.
* Simple, only use standard library besides a logger.
* Graceful shutdown.
* Support for https target servers.

## How to Set Up
#### Prerequisites
* [Golang](https://golang.org/doc/install) >1.19

### Configuration
The app uses a configuration file in TOML format.
The configuration path is set with -config command option when running the app.
The app will use an example configuration file in internal/config/example_config.toml if not provided.

[Toml example](https://github.com/sgrodriguez/reverse-proxy/blob/main/internal/config/example_config.toml)
```toml
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

#### Run the examples
Run the example to test against a mock target server:
```
make example
```

After running make example, you can test the reverse proxy with the target_server_example located in test/ folder.
```
curl -X GET -v http://localhost:8081/get 
```

#### Run the tests
```
make tests
```

## Future Work - Nice To Have
* In order to be production ready it needs more work with:
  * Streaming.
  * Websockets.
  * Compresses data.
  * More testing with the mask to avoid leaks.
* Healthcheck, readiness endpoint.
* Expose the reverse proxy with TLS support.
* Add Timeout configuration to the reverse proxy.
* Metrics.
* Benchmarking and performance testing:
  * Improve the mask regexes with others libs [hyperscan](https://pkg.go.dev/github.com/flier/gohs/hyperscan) [re2](https://github.com/google/re2)
  * Check blockers concurrently.
* Add more test cases of strange and edge cases.

## Contributors
* [sgrodriguez](https://github.com/sgrodriguez)
