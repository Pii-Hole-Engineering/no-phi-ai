#!/usr/bin/env bash
## File:  scripts/make.tidy.sh

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

cd ${_parent_dir} && go mod tidy \
	&& echo "SUCCESS : no-phi-ai : tidy" \
	|| (echo "ERROR : no-phi-ai : tidy" && exit 50)

