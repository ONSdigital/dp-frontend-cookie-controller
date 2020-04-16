dp-frontend-cookie-controller
================

### Getting started

This project uses go modules, ensure that go version 1.12 or above is in use.
If on go 1.12 then ensure the project either resides outside of your `GOPATH` or `GO111MODULE` is set to true

1. To start the service use make, `make debug`

### Dependencies

[dp-frontend-renderer](https://github.com/ONSdigital/dp-frontend-renderer)

### Configuration

| Environment variable          | Config                        | Description                                                                            | Default                 |
| ------------------------------|-------------------------------|----------------------------------------------------------------------------------------| ------------------------|
| BIND_ADDR                     | BindAddr                      | The Port to run on                                                                     | :24100                  |
| RENDERER_URL                  | RendererURL                   | URL dp-frontend-renderer can be reached                                                | https://localhost:20010 |
| SITE_DOMAIN                   | SiteDomain                    |                                                                                        | localhost
| GRACEFUL_SHUTDOWN_TIMEOUT     | GracefulShutdownTimeout       | Time to wait during graceful shutdown                                                  | 5 seconds               |
| HEALTHCHECK_INTERVAL          | HealthCheckInterval           | Interval between health checks                                                         | 30 seconds              |
| HEALTHCHECK_CRITICAL_TIMEOUT  | HealthCheckCriticalTimeout    | Amount of time to pass since last healthy health check to be deemed a critical failure | 90 seconds              |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
