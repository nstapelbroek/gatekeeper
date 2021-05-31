# Gatekeeper

[![Build Status](https://api.cirrus-ci.com/github/nstapelbroek/gatekeeper.svg)](https://cirrus-ci.com/github/nstapelbroek/gatekeeper)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstapelbroek/gatekeeper)](https://goreportcard.com/report/github.com/nstapelbroek/gatekeeper)
[![Image details](https://images.microbadger.com/badges/image/nstapelbroek/gatekeeper.svg)](https://hub.docker.com/r/nstapelbroek/gatekeeper/tags/)
[![Image git ref](https://images.microbadger.com/badges/commit/nstapelbroek/gatekeeper.svg)](https://microbadger.com/images/nstapelbroek/gatekeeper/)
[![GitHub license](https://img.shields.io/github/license/nstapelbroek/gatekeeper.svg)](https://github.com/nstapelbroek/gatekeeper/blob/master/LICENSE)

Temporary allow access to your cloud infrastructure by signaling the gatekeeper. Allowing your build pipeline to deploy behind a firewall.

## Supported environments

| Provider | Product Name | Required Environment Variables |
|---    |---    |---    |
| Vultr | Firewall | `VULTR_PERSONAL_ACCESS_TOKEN`, `VULTR_FIREWALL_ID`|
| Digitalocean | Cloud Firewalls | `DIGITALOCEAN_PERSONAL_ACCESS_TOKEN`, `DIGITALOCEAN_FIREWALL_ID` |
| AWS | EC2 Security Groups | `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `AWS_REGION`, `AWS_SECURITY_GROUP_ID` |
| AWS | VPC Network ACLs | `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `AWS_REGION`, `AWS_NETWORK_ACL_ID` |

## Getting Started

### Installation
1. Download a release binary or use a Docker image
1. Retrieve your cloud provider API keys. [DigitalOcean](https://www.digitalocean.com/docs/api/create-personal-access-token/) even has docs for this.
1. Configure your application by passing environment variables. See these examples below:

Docker:
```
docker run -p 8080:8080 -e DIGITALOCEAN_PERSONAL_ACCESS_TOKEN=REPLACE_ME -e DIGITALOCEAN_FIREWALL_ID=REPLACE_ME nstapelbroek/gatekeeper:latest
```

Standalone binary:
```
DIGITALOCEAN_PERSONAL_ACCESS_TOKEN=REPLACE_ME DIGITALOCEAN_FIREWALL_ID=REPLACE_ME ./gatekeeper
```

### Usage
After installing and running the application you can fire an HTTP POST towards it to temporary whitelist your given IP at the cloud provider.
By default the gatekeeper will open TCP port 22 (for SSH). You can change the port of protocol in the [configuration](#configuration).

A simple example:
```bash
curl -X POST http://localhost:8080
```

You can configure the timeout or ip address per request basis by sending it as a form-encoded or json payload. The example below will use your public IP:
```bash
curl -X POST -s -d 'ip='$(curl -s https://ifconfig.co/ip)'&timeout=60' http://localhost:8080
```

  
### Configuration

Although this tool is meant to be very simple, you can configure it to your needs by changing some variables. 

| Variable Name      | Default value | Notes |
|---	             |---	        |---    |
| APP_ENV            | release      | Used to control the verbosity of log lines. Only `release` and `debug` are used. |
| HTTP_AUTH_USERNAME |              | Used with to `HTTP_AUTH_PASSWORD` to shield the application with http basic auth. |
| HTTP_AUTH_PASSWORD |              | See `HTTP_AUTH_USENAME`. Both values have to be provided.                         |
| HTTP_PORT          | 8080         | Controls on which port the HTTP server will start.                                |
| RULE_CLOSE_TIMEOUT | 120          | When no timeout value is given on a request, this value in seconds will be used. Use 0 to permanently allow the IP address. |
| RULE_PORTS         | TCP:22       | A comma separated list of ports to unblock on a request. Use a `-` to indicate a range. For example: `TCP:20-22,UDP:20-22`. |


### Development
If you wish to help building gatekeeper you can start with:

1. [Fork and clone the repository](https://github.com/nstapelbroek/gatekeeper/fork)
1. Install dependencies with `go mod tidy`
1. Optionally you can install additional tooling like [golangci-lint](https://github.com/golangci/golangci-lint)
1. Start building! You can find some inspiration for changes in the [issues](https://github.com/nstapelbroek/gatekeeper/issue) or [project board](https://github.com/nstapelbroek/gatekeeper/projects)
