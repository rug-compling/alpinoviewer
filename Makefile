.PHONY: help
help:   ## deze uitleg
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[34m%-20s\033[0m %s\n", $$1, $$2}'

build: ## build in huidige dir
	go build .

build_nodbxml: ## build in huidige dir zonder DbXML
	go build -tag nodbxml .

build_wk41: ## build in huidige dir met WebKit 4.1
	go build -tag wk41.

build_nodbxml_wk41: ## build in huidige dir zonder DbXML met WebKit 4.1
	go build -tag nodbxml,wk41 .

install: ## installeren
	go install .

install_nodbxml: ## installeren zonder DbXML
	go install -tags nodbxml .

install_wk41: ## installeren met WebKit 4.1
	go install -tags wk41 .

install_nodbxml_41: ## installeren zonder DbXML met WebKit 4.1
	go install -tags nodbxml,wk41 .
