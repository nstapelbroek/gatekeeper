# Gatekeeper

Temporary allow access to your cloud infrastructure by signaling the gatekeeper. Allowing your build pipeline to deploy behind a firewall.

## Supported environments

| Provider   | Supported status | Required Environment Variables |
|---	|---	|---    |
| Vultr | :white_check_mark: |`VULTR_PERSONAL_ACCESS_TOKEN`, `VULTR_FIREWALL_ID`|
| Digitalocean |  :white_check_mark: | `DIGITALOCEAN_PERSONAL_ACCESS_TOKEN`, `DIGITALOCEAN_FIREWALL_ID` |
| AWS (Security Groups) | :white_check_mark: | `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, `AWS_REGION`, `AWS_FIREWALL_ID` |
| AWS (Network ACLs) | :construction_worker: | None yet... |

## Getting Started

### Installation
1. Download a release binary or use a Docker image
1. Generate your cloud provider API keys. [DigitalOcean](https://www.digitalocean.com/docs/api/create-personal-access-token/) even has docs for this.
1. Configure your application by passing environment variables. See these examples below:

Docker:
```
docker run -p 8080:8080 -e DIGITALOCEAN_PERSONAL_ACCESS_TOKEN=REPLACE_ME -e DIGITALOCEAN_FIREWALL_ID=REPLACE_ME nstapelbroek/gatekeeper:latest
```

Binary:
```
DIGITALOCEAN_PERSONAL_ACCESS_TOKEN=REPLACE_ME DIGITALOCEAN_FIREWALL_ID=REPLACE_ME ./gatekeeper
```

### Usage
After installing and running the application you can fire an HTTP POST towards it to temporary whitelist your given IP at the cloud provider.
By default, the gatekeeper will open TCP port 22 (for SSH).

A curl example that requests the gatekeeper using your public IP:
```curl
curl -X POST -s -d 'ip='$(curl -s https://ifconfig.co/ip)'&timeout=60' http://localhost:8080
```

Note that you do not need to pass an IP or timeout as form-encoded / json data. A simple POST will use a default timeout
and the remote address to apply the rule. So a simple `curl -X POST http://localhost:8080/` will also work :)

  
### Configuration

Although the tool is very simple, you can configure it to your needs by changing some variables. 

| Variable Name      | Default value | Notes |
|---	             |---	        |---    |
| APP_ENV            | release      | Used to control the verbosity of log lines. Only `release` and `debug` are used. |
| HTTP_AUTH_USERNAME |              | Used with to `HTTP_AUTH_PASSWORD` to shield the application with http basic auth. |
| HTTP_AUTH_PASSWORD |              | See `HTTP_AUTH_USENAME`. Both values have to be provided.                         |
| HTTP_PORT          | 8080         | Controls on which port the HTTP server will start.                                |
| RULE_CLOSE_TIMEOUT | 120          | When no timeout value is given on a request, this value in seconds will be used. Use 0 to permanently allow the IP address. |
| RULE_PORTS         | TCP:22       | A comma separated list of ports to unblock on a request. You cannot overwrite this on request basis. Use a `-` to indicate a range. For example: `TCP:20-22,UDP:20-22`. |


### Development
If you wish to help building gatekeeper you can start with:

1. [Forking the repository](https://github.com/nstapelbroek/gatekeeper/fork)
1. Installing it locally (`go get github.com/{your-github-username}/gatekeeper`)
1. Installing some required tooling like [dep](https://github.com/golang/dep) and [golint](https://github.com/golangci/golangci-lint)
1. Installing the dependencies using `dep ensure`
1. Scouting an [issue](https://github.com/nstapelbroek/gatekeeper/issue) or [backlog task](https://github.com/nstapelbroek/gatekeeper/projects) you want to solve
