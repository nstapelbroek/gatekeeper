.PHONY: build, test, run

PROJECTNAME=nstapelbroek/gatekeeper
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))
PROJECT_PACKAGES=$(shell go list ./... | grep -v vendor)

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o gatekeeper .
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull .
	rm gatekeeper

test:
	golint -set_exit_status -min_confidence=0.9 $(PROJECT_PACKAGES)
	go test $(PROJECT_PACKAGES)

run:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker run --rm --name gatekeeper-run -p 8080:8080 -d $(PROJECTNAME):$(TAGNAME_CLEAN)