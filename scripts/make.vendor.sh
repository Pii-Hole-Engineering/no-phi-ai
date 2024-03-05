#!/usr/bin/env bash
## File:  scripts/make.vendor.sh

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

cd ${_parent_dir} && go mod vendor -v \
	&& echo "SUCCESS : no-phi-ai : vendor" \
	|| (echo "ERROR : no-phi-ai : vendor" && exit 20)

