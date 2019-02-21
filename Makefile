TEST ?= $$(go list ./...)
ifndef POSTGRES_DATA_SOURCE
	export POSTGRES_DATA_SOURCE=postgres://postgres@/terraform_provider_sql?sslmode=disable
endif
ifndef MYSQL_DATA_SOURCE
	export MYSQL_DATA_SOURCE=root@/terraform_provider_sql
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
