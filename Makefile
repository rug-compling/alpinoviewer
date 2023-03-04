.PHONY: help
help:   ## deze uitleg
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[34m%-10s\033[0m %s\n", $$1, $$2}'

build: alpinoviewer ## build in huidige dir

alpinoviewer: *.go *.c *.h
	go build .

install: ## installeren
	go install .
