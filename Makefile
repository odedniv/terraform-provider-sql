TEST ?= $$(go list ./...)
ifndef PGCONN
	export PGCONN=postgres://postgres@localhost/terraform_provider_sql?sslmode=disable
endif

default: build
.PHONY: default

help:
	@echo "Main commands:"
	@echo "  help            - show this message"
	@echo "  build (default) - build the terraform provider"
	@echo "  testacc         - runs acceptance tests"
.PHONY: help

build:
	go build
.PHONY: build

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)
.PHONY: testacc
