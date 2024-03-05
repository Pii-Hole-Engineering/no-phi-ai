#!/bin/bash

set -ex

#
# Install Go tools for VSCode
#
go install -v github.com/cweill/gotests/gotests@latest &&\
go install -v github.com/fatih/gomodifytags@latest &&\
go install -v github.com/josharian/impl@latest &&\
go install -v github.com/go-delve/delve/cmd/dlv@latest &&\
go install -v honnef.co/go/tools/cmd/staticcheck@latest &&\
go install -v golang.org/x/tools/gopls@latest

# install other misc tools
sudo apt-get update && sudo apt-get install -y zsh

cd /workspaces/no-phi-ai

# build the project
make build_container
