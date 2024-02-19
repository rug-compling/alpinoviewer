.PHONY: help
help:   ## deze uitleg
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[34m%-15s\033[0m %s\n", $$1, $$2}'

build: ## build in huidige dir
	go build .

build_nodbxml: ## build in huidige dir zonder DbXML
	go build -tag nodbxml .

install: ## installeren
	go install .

install_nodbxml: ## installeren zonder DbXML
	go install -tags nodbxml .
