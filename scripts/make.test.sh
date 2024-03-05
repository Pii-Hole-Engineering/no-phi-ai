#!/usr/bin/env bash
## File:  scripts/make.test.sh

_verbose=""

if [[ "$1" == "-v" || $1 == "--verbose" ]]
then
	_verbose="-v"
fi

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

cd ${_parent_dir} && go list ./pkg/... | grep -v 'vendor' | xargs go test -buildmode='pie' -cover -coverprofile=coverage.out -timeout=30s ${_verbose} \
	&& echo "SUCCESS : no-phi-ai : test" \
	|| (echo "ERROR : no-phi-ai : test" && exit 60)

