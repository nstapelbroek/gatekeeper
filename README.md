# Gatekeeper

Temporary allow access to your cloud infrastructure by authenticating at the gatekeeper. Allowing your CI to deploy behind a firewall.

## Supported adapters

### Environments
| Provider   | Supported status | Required Environment Variables |
|---	|---	|---    |
| Vultr | Partially supported, no port-ranges and the white list lasts for a hardcoded amount of 120 seconds |`VULTR_API_KEY`, `VULTR_FIREWALL_GROUP_ID`|
| Digitalocean | In development | None yet |

## Getting Started

No instructions yet  ¯\_(ツ)_/¯

### For development
If you wish to help building gatekeeper you can start with:

1. [Forking the repository](https://github.com/nstapelbroek/gatekeeper/fork)
1. Installing it locally (`go get github.com/{your-github-username}/gatekeeper`)
1. Installing some required tooling like [dep](https://github.com/golang/dep) and [golint](https://github.com/golang/lint)
1. Installing the dependencies using `dep ensure`
1. Scouting an [issue](https://github.com/nstapelbroek/gatekeeper/issue) or [backlog task](https://github.com/nstapelbroek/gatekeeper/projects) you want to solve