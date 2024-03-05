#!/usr/bin/env bash
## File:  scripts/make.build_only.sh

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

cd ${_parent_dir}/pkg/ \
	&& go build -buildmode="pie" -buildvcs=false -o ${_parent_dir}/build/no-phi-ai \
	&& echo "SUCCESS : no-phi-ai : build_only" \
	|| (echo "ERROR : no-phi-ai : build_only" && exit 70)

