# Makefile : no-phi-ai
#
_app_name=no-phi-ai
_app_pkg_dir=pkg
_build_dir=build
_coverage_out_file=coverage.out

_cmd_docker_build=DOCKER_BUILDKIT=1 docker build --ssh default
_cmd_go_cover=go tool cover -func=${_coverage_out_file}

_msg="${_app_name} : make"
_msg_error="ERROR : ${_msg}"
_msg_success="SUCCESS : ${_msg}"

default: build

.PHONY: build clean cover deploy format image package remove test tidy vendor

build: build_prep build_only

build_container: vendor build_only

build_only:
	./scripts/make.build_only.sh

build_prep: clean format tidy test

clean: clean_test
	rm -rf ./${_build_dir}/${_app_name} ./vendor/ Gopkg.lock ${_coverage_out_file}

clean_and_build: clean build

clean_test:
	go clean -testcache

cover: build_only test cover_only

cover_only:
	${_cmd_go_cover}

format:
	./scripts/make.format.sh

image: image_mini

image_full:
	./scripts/make.image.sh "full"

image_mini:
	./scripts/make.image.sh

test:
	./scripts/make.test.sh

test_build: tidy build_only test

test_full: clean_test test_verbose cover

test_verbose:
	./scripts/make.test.sh --verbose

tidy:
	./scripts/make.tidy.sh

vendor:
	./scripts/make.vendor.sh
