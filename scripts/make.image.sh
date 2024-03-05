#!/usr/bin/env bash
## File:  scripts/make.image.sh

_exit=80
_tag="no-phi-ai"
_target="mini"

if [[ "$1" == "full" ]]
then
	_exit=85
	_tag="${_tag}-full"
	_target="full"
fi

_script_dir=`cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd`
_parent_dir=${_script_dir}/..

export DOCKER_BUILDKIT=1

cd ${_parent_dir} && docker build --ssh default --tag ${_tag} --target ${_target} . \
	&& echo "SUCCESS : no-phi-ai : image=${_tag}" \
	|| (echo "ERROR : no-phi-ai : image=${_tag}" && exit ${_exit})

