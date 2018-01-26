.PHONY: build

PROJECTNAME=nstapelbroek/gatekeeper
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o gatekeeper .
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull .
	rm gatekeeper

run:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker run --rm --name gatekeeper-run -p 8888:8888 -d $(PROJECTNAME):$(TAGNAME_CLEAN)