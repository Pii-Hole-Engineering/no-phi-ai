#!/usr/bin/env bash
## File:  scripts/make.format.sh

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

cd ${_parent_dir} && go fmt ./pkg/... \
	&& echo "SUCCESS : no-phi-ai : format" \
	|| (echo "ERROR : no-phi-ai : format" && exit 40)

