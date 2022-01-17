.PHONY: build, test, run

PROJECTNAME=nstapelbroek/gatekeeper
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))
GIT_REV=$(shell git rev-parse HEAD)

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull --build-arg VCS_REF=$(GIT_REV) .

test:
	golangci-lint run
	go test ./...

publish:
	KO_DOCKER_REPO=$(PROJECTNAME) ko publish ./cmd/gatekeeper --platform=linux/amd64,linux/arm,linux/arm64 --bare

run:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker run --rm --name gatekeeper-run -p 8080:8080 -e VULTR_PERSONAL_ACCESS_TOKEN=somekey -e VULTR_FIREWALL_ID=somegroup -d $(PROJECTNAME):$(TAGNAME_CLEAN)
