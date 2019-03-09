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

run:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker run --rm --name gatekeeper-run -p 8080:8080 -e VULTR_API_KEY=somekey -e VULTR_FIREWALL_GROUP=somegroup -d $(PROJECTNAME):$(TAGNAME_CLEAN)
